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
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bobg/mid"
	"github.com/bobg/sqlutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/bobg/tredd"
)

func serve(args []string) {
	ctx := context.Background()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		addr   = fs.String("addr", "localhost:20544", "server listen address")
		dir    = fs.String("dir", ".", "root of content tree")
		dbFile = fs.String("db", "", "file containing server-state db")
		ethURL = fs.String("ethurl", "", "URL of blockchain server")
	)

	keyfile, passphrase := addKeyfilePassphrase(fs)

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	db, err := openDB(ctx, *dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seller, err := handleKeyfilePassphrase(*keyfile, *passphrase)
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		db:     db,
		dir:    *dir,
		seller: seller,
		ethURL: *ethURL,
	}

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

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", listener.Addr())

	http.Handle("/request", mid.Err(s.serve))
	// http.HandleFunc("/propose-payment", s.revealKey)
	http.Serve(listener, nil)
}

type server struct {
	db     *sql.DB // transfer records and blockchain info
	dir    string  // content
	seller *bind.TransactOpts
	ethURL string
}

type serverRecord struct {
	tredd.ParseResult
	transferID [32]byte
}

const (
	minRevealDur = 10 * time.Minute
	maxRefundDur = time.Hour
)

func (s *server) serve(w http.ResponseWriter, req *http.Request) error {
	var (
		clearRootStr      = req.FormValue("clearroot")
		amountStr         = req.FormValue("amount")
		tokenType         = req.FormValue("token")
		revealDeadlineStr = req.FormValue("revealdeadline")
		refundDeadlineStr = req.FormValue("refunddeadline")
	)

	var clearRoot [32]byte
	_, err := hex.Decode(clearRoot[:], []byte(clearRootStr))
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "decoding clear root")}
	}

	dir, filename := clearHashPath(s.dir, clearRoot)
	f, err := os.Open(path.Join(dir, filename))
	if os.IsNotExist(err) {
		return mid.CodeErr{C: http.StatusNotFound}
	}
	if err != nil {
		return errors.Wrapf(err, "opening %s", filename)
	}
	defer f.Close()

	contentType, err := ioutil.ReadFile(path.Join(dir, "content-type"))
	if err != nil {
		return errors.Wrap(err, "getting content type")
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "parsing amount")}
	}
	if amount < 1 {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrapf(err, "non-positive amount %d", amount)}
	}

	err = s.checkPrice(amount, tokenType, clearRoot)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "proposed payment rejected")}
	}

	revealDeadlineMS, err := strconv.ParseUint(revealDeadlineStr, 10, 64)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "parsing reveal deadline")}

	}
	revealDeadline := FromMillis(revealDeadlineMS)

	if time.Until(revealDeadline) < minRevealDur {
		return mid.CodeErr{C: http.StatusBadRequest, Err: fmt.Errorf("reveal deadline too soon: %s, require %s", time.Until(revealDeadline), minRevealDur)}
	}

	refundDeadlineMS, err := strconv.ParseUint(refundDeadlineStr, 10, 64)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "parsing refund deadline")}
	}
	refundDeadline := FromMillis(refundDeadlineMS)

	if refundDeadline.Sub(revealDeadline) > maxRefundDur {
		return mid.CodeErr{C: http.StatusBadRequest, Err: fmt.Errorf("refund deadline too much later after reveal deadline: %s, require %s", refundDeadline.Sub(revealDeadline), maxRefundDur)}
	}

	var key [32]byte
	_, err = rand.Read(key[:])
	if err != nil {
		return errors.Wrap(err, "choosing cipher key")
	}

	rec := &serverRecord{
		ParseResult: tredd.ParseResult{
			Amount:         amount,
			TokenType:      tokenType,
			ClearRoot:      clearRoot,
			RevealDeadline: revealDeadline,
			RefundDeadline: refundDeadline,
			Seller:         s.seller.From, // TODO: check this is right
			Key:            key,
		},
	}

	_, err = rand.Read(rec.transferID[:])
	if err != nil {
		return errors.Wrap(err, "choosing transfer ID")
	}

	log.Printf("new transfer %x, clearRoot %x, payment %d/%s, deadlines %s/%s, key %x", rec.transferID[:], clearRoot, amount, tokenType, revealDeadline, refundDeadline, key[:])

	w.Header().Set("X-Tredd-Transfer-Id", hex.EncodeToString(rec.transferID[:]))
	w.Header().Set("Content-Type", string(contentType))

	tmpfile, err := ioutil.TempFile("", "treddserve")
	if err != nil {
		return errors.Wrap(err, "creating response tempfile")
	}
	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)
	defer tmpfile.Close()

	cipherRoot, err := tredd.Serve(tmpfile, f, key)
	if err != nil {
		return errors.Wrap(err, "serving data")
	}

	err = tmpfile.Close()
	if err != nil {
		return errors.Wrap(err, "closing response tempfile")
	}

	copy(rec.CipherRoot[:], cipherRoot)

	err = s.storeRecord(req.Context(), rec)
	if err != nil {
		return errors.Wrap(err, "storing transfer record")
	}

	tmpfile, err = os.Open(tmpfilename)
	if err != nil {
		return errors.Wrap(err, "reopening response tempfile")
	}
	defer tmpfile.Close()
	_, err = io.Copy(w, tmpfile)
	if err != nil {
		return errors.Wrap(err, "writing response")
	}

	// TODO: queue a blockchain watcher that does revealKey when a Tredd contract with this cipher root shows up

	return nil
}

// TODO: this is no longer an HTTP entrypoint; it's a callback based on a blockchain event
func (s *server) revealKey(w http.ResponseWriter, req *http.Request) {
	transferIDStr := req.Header.Get("X-Tredd-Transfer-Id")

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

	client, err := ethclient.Dial(s.ethURL)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "contacting Ethereum server: %s", err)
		return
	}

	// TODO: populate these
	var (
		contractAddr, wantTokenType common.Address
		wantAmount, wantCollateral  *big.Int
	)

	receipt, err := tredd.RevealKey(
		ctx,
		client,
		s.seller,
		contractAddr,
		rec.Key,
		wantTokenType,
		wantAmount, wantCollateral,
		rec.ClearRoot, rec.CipherRoot,
	)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "constructing reveal-key transaction: %s", err)
		return
	}

	log.Printf("revealed key in transaction %x", receipt.TxHash[:])

	// TODO: parse the buyer out of the contract, if in fact we need it for storeRecord
	// rec.Buyer = parsed.Buyer

	err = s.storeRecord(ctx, rec)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "updating transfer record")
		return
	}

	s.queueClaimPaymentHelper(ctx, rec)

	log.Printf("transfer %x: revealing key", transferID)

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getRecord(ctx context.Context, transferID []byte) (*serverRecord, error) {
	var rec serverRecord
	copy(rec.transferID[:], transferID)

	const q = `
		SELECT key, contract_addr, clear_root, cipher_root, token_type, amount, reveal_deadline_ms, refund_deadline_ms, buyer, seller
			FROM transfer_records
			WHERE transfer_id = $1
	`

	var (
		revealDeadlineMS, refundDeadlineMS uint64
	)
	err := s.db.QueryRowContext(ctx, q, transferID).Scan(&rec.Key, &rec.ContractAddr, &rec.ClearRoot, &rec.CipherRoot, &rec.TokenType, &rec.Amount, &revealDeadlineMS, &refundDeadlineMS, &rec.Buyer, &rec.Seller)
	if err != nil {
		return nil, errors.Wrapf(err, "querying transfer record %x from db", transferID)
	}
	rec.RevealDeadline = FromMillis(revealDeadlineMS)
	rec.RefundDeadline = FromMillis(refundDeadlineMS)
	return &rec, nil
}

func (s *server) storeRecord(ctx context.Context, rec *serverRecord) error {
	const q = `
		INSERT OR REPLACE INTO transfer_records
			(transfer_id, key, contract_addr, clear_root, cipher_root, token_type, amount, reveal_deadline_ms, refund_deadline_ms, buyer, seller)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := s.db.ExecContext(ctx, q, rec.transferID[:], rec.Key, rec.ContractAddr, rec.ClearRoot, rec.CipherRoot, rec.TokenType, rec.Amount, Millis(rec.RevealDeadline), Millis(rec.RefundDeadline), rec.Buyer, rec.Seller)
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
	// TODO: set a timer that does tredd.ClaimPayment after the refund deadline
	// It should also delete the row from transfer_records.
}

func (s *server) checkPrice(amount int64, tokenType string, clearRoot [32]byte) error {
	if amount > 0 { // TODO: per-content pricing!
		return nil
	}
	return errors.New("amount must be 1 or higher")
}

func httpErrf(w http.ResponseWriter, code int, msgfmt string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(msgfmt, args...), code)
	log.Printf(msgfmt, args...)
}
