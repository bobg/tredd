package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bobg/sqlutil"
	"github.com/bobg/tredd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txvm"
)

func serve(args []string) {
	ctx := context.Background()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		addr    = fs.String("addr", "localhost:20544", "server listen address")
		dir     = fs.String("dir", ".", "root of content tree")
		dbFile  = fs.String("db", "", "file containing server-state db")
		prvFile = fs.String("prv", "", "file containing server private key")
		url     = fs.String("url", "", "URL of blockchain server")
	)

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	submitURL := *url + "/submit"
	getURL := *url + "/get"

	f, err := os.Open(*prvFile)
	if err != nil {
		log.Fatalf("opening prv file %s: %s", *prvFile, err)
	}
	defer f.Close()

	var prvbuf [ed25519.PrivateKeySize]byte
	_, err = io.ReadFull(f, prvbuf[:])
	if err != nil {
		log.Fatalf("reading prv file %s: %s", *prvFile, err)
	}
	f.Close()

	prv := ed25519.PrivateKey(prvbuf[:])

	db, err := openDB(ctx, *dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seller := prv.Public().(ed25519.PublicKey)
	s := &server{
		db:     db,
		dir:    *dir,
		seller: seller,
		o:      newObserver(db, seller, getURL),
	}
	s.signer = func(msg []byte) ([]byte, error) {
		return ed25519.Sign(prv, msg), nil
	}
	s.submitter = submitter(submitURL)

	var transferIDs [][]byte
	err = sqlutil.ForQueryRows(ctx, db, "SELECT transfer_id FROM transfer_records", func(transferID []byte) {
		transferIDs = append(transferIDs, transferID)
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, transferID := range transferIDs {
		log.Printf("queueing claim-payment callback for transfer %x", transferID)
		err = s.queueClaimPayment(ctx, transferID)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("starting blockchain observer")
	go s.o.run(ctx)

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", listener.Addr())

	http.HandleFunc("/request", s.serve)
	http.HandleFunc("/propose-payment", s.revealKey)
	http.Serve(listener, nil)
}

type server struct {
	db        *sql.DB // transfer records and blockchain info
	dir       string  // content
	seller    ed25519.PublicKey
	o         *observer
	signer    tredd.Signer
	submitter func(prog []byte, version, runlimit int64) error
}

type serverRecord struct {
	tredd.ParseResult
	transferID [32]byte
}

const (
	minRevealDur = 10 * time.Minute
	maxRefundDur = time.Hour
)

func (s *server) serve(w http.ResponseWriter, req *http.Request) {
	var (
		clearRootStr      = req.FormValue("clearroot")
		amountStr         = req.FormValue("amount")
		assetIDStr        = req.FormValue("assetid")
		revealDeadlineStr = req.FormValue("revealdeadline")
		refundDeadlineStr = req.FormValue("refunddeadline")
	)

	clearRoot, err := hex.DecodeString(clearRootStr)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "decoding clear root: %s", err)
		return
	}

	dir, filename := clearHashPath(s.dir, clearRoot)
	f, err := os.Open(path.Join(dir, filename))
	if os.IsNotExist(err) {
		httpErrf(w, http.StatusNotFound, "file not found")
		return
	}
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "opening %s: %s", filename, err)
		return
	}
	defer f.Close()

	contentType, err := ioutil.ReadFile(path.Join(dir, "content-type"))
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "getting content type: %s", err)
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "parsing amount: %s", err)
		return
	}
	if amount < 1 {
		httpErrf(w, http.StatusBadRequest, "non-positive amount %d", amount)
		return
	}
	assetID, err := hex.DecodeString(assetIDStr)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "parsing asset ID: %s", err)
		return
	}

	err = s.checkPrice(amount, assetID, clearRoot)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "proposed payment rejected: %s", err)
		return
	}

	revealDeadlineMS, err := strconv.ParseUint(revealDeadlineStr, 10, 64)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "parsing reveal deadline: %s", err)
		return
	}
	revealDeadline := bc.FromMillis(revealDeadlineMS)

	if time.Until(revealDeadline) < minRevealDur {
		httpErrf(w, http.StatusBadRequest, "reveal deadline too soon: %s, require %s", time.Until(revealDeadline), minRevealDur)
		return
	}

	refundDeadlineMS, err := strconv.ParseUint(refundDeadlineStr, 10, 64)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "parsing refund deadline: %s", err)
		return
	}
	refundDeadline := bc.FromMillis(refundDeadlineMS)

	if refundDeadline.Sub(revealDeadline) > maxRefundDur {
		httpErrf(w, http.StatusBadRequest, "refund deadline too later after reveal deadline: %s, require %s", refundDeadline.Sub(revealDeadline), maxRefundDur)
		return
	}

	var key [32]byte
	_, err = rand.Read(key[:])
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "choosing cipher key: %s", err)
		return
	}

	rec := &serverRecord{
		ParseResult: tredd.ParseResult{
			Amount:         amount,
			AssetID:        assetID,
			ClearRoot:      clearRoot,
			RevealDeadline: revealDeadline,
			RefundDeadline: refundDeadline,
			Seller:         s.seller,
			Key:            key[:],
		},
	}

	_, err = rand.Read(rec.transferID[:])
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "choosing transfer ID: %s", err)
		return
	}

	log.Printf("new transfer %x, clearRoot %x, payment %d/%x, deadlines %s/%s, key %x", rec.transferID[:], clearRoot, amount, assetID, revealDeadline, refundDeadline, key[:])

	w.Header().Set("X-Tredd-Transfer-Id", hex.EncodeToString(rec.transferID[:]))
	w.Header().Set("Content-Type", string(contentType))

	tmpfile, err := ioutil.TempFile("", "treddserve")
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "creating response tempfile: %s", err)
		return
	}
	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)
	defer tmpfile.Close()

	cipherRoot, err := tredd.Serve(tmpfile, f, key)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "serving data: %s", err)
		return
	}

	err = tmpfile.Close()
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "closing response tempfile: %s", err)
		return
	}

	rec.CipherRoot = cipherRoot

	err = s.storeRecord(req.Context(), rec)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "storing transfer record: %s", err)
		return
	}

	tmpfile, err = os.Open(tmpfilename)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "reopening response tempfile: %s", err)
		return
	}
	defer tmpfile.Close()
	_, err = io.Copy(w, tmpfile)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "writing response: %s", err)
		return
	}
}

func (s *server) revealKey(w http.ResponseWriter, req *http.Request) {
	transferIDStr := req.Header.Get("X-Tredd-Transfer-Id")

	paymentProposal, err := ioutil.ReadAll(req.Body)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "reading payment proposal: %s", err)
		return
	}

	transferID, err := hex.DecodeString(transferIDStr)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "decoding transfer ID: %s", err)
		return
	}

	ctx := req.Context()
	rec, err := s.getRecord(ctx, transferID)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "finding transfer record: %s", err)
		return
	}

	var (
		clearRoot  [32]byte
		cipherRoot [32]byte
		key        [32]byte
		assetID    = bc.HashFromBytes(rec.AssetID)
	)
	copy(clearRoot[:], rec.ClearRoot)
	copy(cipherRoot[:], rec.CipherRoot)
	copy(key[:], rec.Key)

	now := time.Now()

	prog, err := tredd.RevealKey(ctx, paymentProposal, s.seller, key, rec.Amount, assetID, s.o.r, s.signer, clearRoot, cipherRoot, now, rec.RevealDeadline, rec.RefundDeadline)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "constructing reveal-key transaction: %s", err)
		return
	}

	parsed := tredd.ParseLog(prog)
	if parsed == nil {
		httpErrf(w, http.StatusInternalServerError, "parsing tx log")
		return
	}

	rec.Anchor1 = parsed.Anchor1
	rec.Anchor2 = parsed.Anchor2
	rec.Buyer = parsed.Buyer
	rec.OutputID = parsed.OutputID

	err = s.storeRecord(ctx, rec)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "updating transfer record")
		return
	}

	vm, err := txvm.Validate(prog, 3, math.MaxInt64)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "computing runlimit: %s", err)
		return
	}

	s.queueClaimPaymentHelper(ctx, rec)

	log.Printf("transfer %x: revealing key", transferID)

	err = s.submitter(prog, 3, math.MaxInt64-vm.Runlimit())
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "submitting reveal-key transaction: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getRecord(ctx context.Context, transferID []byte) (*serverRecord, error) {
	var rec serverRecord
	copy(rec.transferID[:], transferID)

	const q = `
		SELECT key, output_id, clear_root, cipher_root, asset_id, amount, anchor1, anchor2, reveal_deadline_ms, refund_deadline_ms, buyer, seller
			FROM transfer_records
			WHERE transfer_id = $1
	`

	var (
		revealDeadlineMS, refundDeadlineMS uint64
		buyer, seller                      []byte
	)
	err := s.db.QueryRowContext(ctx, q, transferID).Scan(&rec.Key, &rec.OutputID, &rec.ClearRoot, &rec.CipherRoot, &rec.AssetID, &rec.Amount, &rec.Anchor1, &rec.Anchor2, &revealDeadlineMS, &refundDeadlineMS, &buyer, &seller)
	if err != nil {
		return nil, errors.Wrapf(err, "querying transfer record %x from db", transferID)
	}
	rec.Buyer = ed25519.PublicKey(buyer)
	rec.Seller = ed25519.PublicKey(seller)
	rec.RevealDeadline = bc.FromMillis(revealDeadlineMS)
	rec.RefundDeadline = bc.FromMillis(refundDeadlineMS)
	return &rec, nil
}

func (s *server) storeRecord(ctx context.Context, rec *serverRecord) error {
	const q = `
		INSERT OR REPLACE INTO transfer_records
			(transfer_id, key, output_id, clear_root, cipher_root, asset_id, amount, anchor1, anchor2, reveal_deadline_ms, refund_deadline_ms, buyer, seller)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := s.db.ExecContext(ctx, q, rec.transferID[:], rec.Key, rec.OutputID, rec.ClearRoot, rec.CipherRoot, rec.AssetID, rec.Amount, rec.Anchor1, rec.Anchor2, bc.Millis(rec.RevealDeadline), bc.Millis(rec.RefundDeadline), []byte(rec.Buyer), []byte(rec.Seller))
	return err
}

func (s *server) queueClaimPayment(ctx context.Context, transferID []byte) error {
	rec, err := s.getRecord(ctx, transferID)
	if err != nil {
		return err
	}
	s.queueClaimPaymentHelper(ctx, rec)
	return nil
}

func (s *server) queueClaimPaymentHelper(ctx context.Context, rec *serverRecord) {
	s.o.enqueue(rec.RefundDeadline, func() {
		redeem := &tredd.Redeem{
			RefundDeadline: rec.RefundDeadline,
			Buyer:          rec.Buyer,
			Seller:         s.seller,
			Amount:         rec.Amount,
			AssetID:        bc.HashFromBytes(rec.AssetID),
		}
		copy(redeem.Anchor2[:], rec.Anchor2)
		copy(redeem.CipherRoot[:], rec.CipherRoot)
		copy(redeem.ClearRoot[:], rec.ClearRoot)
		copy(redeem.Key[:], rec.Key)

		prog, err := tredd.ClaimPayment(redeem)
		if err != nil {
			log.Fatalf("constructing claim-payment transaction: %s", err)
		}
		vm, err := txvm.Validate(prog, 3, math.MaxInt64)
		if err != nil {
			log.Fatalf("computing runlimit for claim-payment transaction: %s", err)
		}
		err = s.submitter(prog, 3, math.MaxInt64-vm.Runlimit())
		if err != nil {
			log.Fatalf("submitting claim-payment transaction: %s", err) // xxx this one should prob have a retry loop
		}
		_, err = s.db.ExecContext(ctx, "DELETE FROM transfer_records WHERE transfer_id = $1", rec.transferID[:])
		if err != nil {
			log.Printf("WARNING: could not delete transfer record %x: %s", rec.transferID[:], err)
		}
	})
}

func (s *server) checkPrice(amount int64, assetID []byte, clearRoot []byte) error {
	if amount > 0 { // TODO: per-content pricing!
		return nil
	}
	return errors.New("amount must be 1 or higher")
}

func httpErrf(w http.ResponseWriter, code int, msgfmt string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(msgfmt, args...), code)
	log.Printf(msgfmt, args...)
}
