package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bobg/tedd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txvm"
	"github.com/coreos/bbolt"
)

func serve(args []string) {
	ctx := context.Background()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		listen  = fs.String("listen", "", "listen address")
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

	db, err := bbolt.Open(*dbFile, 0600, nil)
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
	err = db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		if root == nil {
			return nil
		}
		recordsBucket := root.Bucket([]byte("records"))
		if recordsBucket == nil {
			return nil
		}
		return recordsBucket.ForEach(func(transferID, _ []byte) error {
			transferIDs = append(transferIDs, transferID)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, transferID := range transferIDs {
		err = s.queueClaimPayment(transferID)
		if err != nil {
			log.Fatal(err)
		}
	}

	go s.o.run(ctx)

	http.HandleFunc("/request", s.serve)
	http.HandleFunc("/propose-payment", s.revealKey)
	http.ListenAndServe(*listen, nil)
}

type server struct {
	db        *bbolt.DB // transfer records and blockchain info
	dir       string    // content
	seller    ed25519.PublicKey
	o         *observer
	signer    tedd.Signer
	submitter func(prog []byte, version, runlimit int64) error
}

type serverRecord struct {
	tedd.ParseResult
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

	// xxx check amount/assetID is acceptable for clearRoot

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
		ParseResult: tedd.ParseResult{
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

	w.Header().Set("X-Tedd-Transfer-Id", hex.EncodeToString(rec.transferID[:]))
	w.Header().Set("Content-Type", string(contentType))

	tmpfile, err := ioutil.TempFile("", "teddserve")
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "creating response tempfile: %s", err)
		return
	}
	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)
	defer tmpfile.Close()

	cipherRoot, err := tedd.Serve(tmpfile, f, key)
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

	err = s.storeRecord(rec)
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
	var (
		transferIDStr      = req.FormValue("transferid")
		paymentProposalStr = req.FormValue("paymentproposal")
	)

	transferID, err := hex.DecodeString(transferIDStr)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "decoding transfer ID: %s", err)
		return
	}

	rec, err := s.getRecord(transferID)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "finding transfer record: %s", err)
		return
	}

	paymentProposal, err := hex.DecodeString(paymentProposalStr)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "decoding payment proposal: %s", err)
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

	now, err := s.o.now()
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "getting blockchain time: %s", err)
		return
	}

	prog, err := tedd.RevealKey(req.Context(), paymentProposal, s.seller, key, rec.Amount, assetID, s.o.r, s.signer, clearRoot, cipherRoot, now, rec.RevealDeadline, rec.RefundDeadline)
	if err != nil {
		httpErrf(w, http.StatusBadRequest, "constructing reveal-key transaction: %s", err)
		return
	}

	parsed := tedd.ParseLog(prog)
	if parsed == nil {
		httpErrf(w, http.StatusInternalServerError, "parsing tx log")
		return
	}

	rec.Anchor1 = parsed.Anchor1
	rec.Anchor2 = parsed.Anchor2
	rec.Buyer = parsed.Buyer
	rec.OutputID = parsed.OutputID

	err = s.storeRecord(rec)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "updating transfer record")
		return
	}

	vm, err := txvm.Validate(prog, 3, math.MaxInt64)
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "computing runlimit: %s", err)
		return
	}

	s.queueClaimPaymentHelper(rec)

	err = s.submitter(prog, 3, math.MaxInt64-vm.Runlimit())
	if err != nil {
		httpErrf(w, http.StatusInternalServerError, "submitting reveal-key transaction: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getRecord(transferID []byte) (*serverRecord, error) {
	var rec serverRecord
	copy(rec.transferID[:], transferID)
	err := s.db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("root"))
		if root == nil {
			return errors.New("no root bucket")
		}
		recordsBucket := root.Bucket([]byte("records"))
		if recordsBucket == nil {
			return errors.New("no records bucket")
		}
		bu := recordsBucket.Bucket(transferID)
		if bu == nil {
			return fmt.Errorf("no record bucket %x", transferID)
		}
		copy(rec.Key[:], bu.Get([]byte("key")))
		rec.ClearRoot = bu.Get([]byte("clearRoot"))
		rec.CipherRoot = bu.Get([]byte("cipherRoot"))
		rec.AssetID = bu.Get([]byte("assetID"))

		var n int
		rec.Amount, n = binary.Varint(bu.Get([]byte("amount")))
		if n < 1 {
			return fmt.Errorf("cannot parse amount in record %x", transferID)
		}
		revealDeadlineMS, n := binary.Uvarint(bu.Get([]byte("revealDeadlineMS")))
		if n < 1 {
			return fmt.Errorf("cannot parse reveal deadline in record %x", transferID)
		}
		rec.RevealDeadline = bc.FromMillis(revealDeadlineMS)
		refundDeadlineMS, n := binary.Uvarint(bu.Get([]byte("refundDeadlineMS")))
		if n < 1 {
			return fmt.Errorf("cannot parse refund deadline in record %x", transferID)
		}
		rec.RefundDeadline = bc.FromMillis(refundDeadlineMS)
		return nil
	})
	return &rec, err
}

func (s *server) storeRecord(rec *serverRecord) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return errors.Wrap(err, "getting/creating root bucket")
		}
		records, err := root.CreateBucketIfNotExists([]byte("records"))
		if err != nil {
			return errors.Wrap(err, "getting/creating records bucket")
		}
		bu, err := records.CreateBucketIfNotExists(rec.transferID[:])
		if err != nil {
			return errors.Wrapf(err, "creating record bucket %x", rec.transferID[:])
		}

		var amountBuf [binary.MaxVarintLen64]byte
		m := binary.PutVarint(amountBuf[:], rec.Amount)
		err = bu.Put([]byte("amount"), amountBuf[:m])
		if err != nil {
			return errors.Wrap(err, "storing amount")
		}

		err = bu.Put([]byte("assetID"), rec.AssetID)
		if err != nil {
			return errors.Wrap(err, "storing assetID")
		}

		err = bu.Put([]byte("anchor1"), rec.Anchor1)
		if err != nil {
			return errors.Wrap(err, "storing anchor1")
		}

		err = bu.Put([]byte("anchor2"), rec.Anchor2)
		if err != nil {
			return errors.Wrap(err, "storing anchor2")
		}

		err = bu.Put([]byte("clearRoot"), rec.ClearRoot)
		if err != nil {
			return errors.Wrap(err, "storing clearRoot")
		}

		err = bu.Put([]byte("cipherRoot"), rec.CipherRoot)
		if err != nil {
			return errors.Wrap(err, "storing cipherRoot")
		}

		var revealDeadlineMSBuf [binary.MaxVarintLen64]byte
		m = binary.PutUvarint(revealDeadlineMSBuf[:], bc.Millis(rec.RevealDeadline))
		err = bu.Put([]byte("revealDeadlineMS"), revealDeadlineMSBuf[:m])
		if err != nil {
			return errors.Wrap(err, "storing reveal deadline")
		}

		var refundDeadlineMSBuf [binary.MaxVarintLen64]byte
		m = binary.PutUvarint(refundDeadlineMSBuf[:], bc.Millis(rec.RefundDeadline))
		err = bu.Put([]byte("refundDeadlineMS"), refundDeadlineMSBuf[:m])
		if err != nil {
			return errors.Wrap(err, "storing refund deadline")
		}

		err = bu.Put([]byte("buyer"), rec.Buyer)
		if err != nil {
			return errors.Wrap(err, "storing buyer")
		}

		err = bu.Put([]byte("seller"), rec.Seller)
		if err != nil {
			return errors.Wrap(err, "storing seller")
		}

		err = bu.Put([]byte("key"), rec.Key)
		if err != nil {
			return errors.Wrap(err, "storing key")
		}

		err = bu.Put([]byte("outputID"), rec.OutputID)
		if err != nil {
			return errors.Wrap(err, "storing outputID")
		}

		return nil
	})
}

func (s *server) queueClaimPayment(transferID []byte) error {
	rec, err := s.getRecord(transferID)
	if err != nil {
		return err
	}
	s.queueClaimPaymentHelper(rec)
	return nil
}

func (s *server) queueClaimPaymentHelper(rec *serverRecord) {
	s.o.enqueue(rec.RefundDeadline, func() {
		redeem := &tedd.Redeem{
			RefundDeadline: rec.RefundDeadline,
			Buyer:          rec.Buyer,
			Seller:         s.seller,
			Amount:         rec.Amount,
			AssetID:        bc.HashFromBytes(rec.AssetID),
		}
		copy(redeem.Anchor[:], rec.Anchor2)
		copy(redeem.CipherRoot[:], rec.CipherRoot)
		copy(redeem.ClearRoot[:], rec.ClearRoot)
		copy(redeem.Key[:], rec.Key)

		prog, err := tedd.ClaimPayment(redeem)
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
		err = s.db.Update(func(tx *bbolt.Tx) error {
			root := tx.Bucket([]byte("root"))         // xxx check
			records := root.Bucket([]byte("records")) // xxx check
			return records.DeleteBucket(rec.transferID[:])
		})
	})
}

func httpErrf(w http.ResponseWriter, code int, msgfmt string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(msgfmt, args...), code)
	log.Printf(msgfmt, args...)
}
