package main

import (
	"bytes"
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
	"github.com/golang/protobuf/proto"
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
		reserver: &reserver{
			pubkey: seller,
			db:     db,
		},
	}
	s.signer = func(msg []byte) ([]byte, error) {
		return ed25519.Sign(prv, msg), nil
	}
	s.submitter = func(prog []byte, version, runlimit int64) error {
		rawTx := &bc.RawTx{
			Version:  version,
			Runlimit: runlimit,
			Program:  prog,
		}
		bits, err := proto.Marshal(rawTx)
		if err != nil {
			return errors.Wrap(err, "serializing tx")
		}
		resp, err := http.Post(submitURL, "application/octet-stream", bytes.NewReader(bits))
		if err != nil {
			return errors.Wrap(err, "submitting tx")
		}
		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("status code %d when submitting tx", resp.StatusCode)
		}
		return nil
	}

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

	go s.monitorBlockchain(ctx, getURL)

	http.HandleFunc("/", s.serve)
	http.HandleFunc("/propose-payment", s.revealKey)
	http.ListenAndServe(*listen, nil)
}

type server struct {
	db        *bbolt.DB // transfer records
	dir       string    // content
	seller    ed25519.PublicKey
	now       time.Time // timestamp of latest blockchain block
	reserver  *reserver // must satisfy tedd.Reserver
	signer    tedd.Signer
	submitter func(prog []byte, version, runlimit int64) error
}

type serverRecord struct {
	transferID            [32]byte
	key                   [32]byte
	clearRoot, cipherRoot []byte
	amount                int64
	assetID               []byte
	revealDeadline        time.Time
	refundDeadline        time.Time
}

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
		// xxx
	}

	dir, filename := clearHashPath(s.dir, clearRoot)
	f, err := os.Open(path.Join(dir, filename))
	if os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("opening %s: %s", filename, err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	contentType, err := ioutil.ReadFile(path.Join(dir, "content-type"))
	if err != nil {
		http.Error(w, fmt.Sprintf("getting content type: %s", err), http.StatusInternalServerError)
		return
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		// xxx
	}
	assetID, err := hex.DecodeString(assetIDStr)
	if err != nil {
		// xxx
	}

	// xxx check amount/assetID is acceptable for clearRoot

	revealDeadlineMS, err := strconv.ParseUint(revealDeadlineStr, 10, 64)
	if err != nil {
		// xxx
	}
	revealDeadline := bc.FromMillis(revealDeadlineMS)

	// xxx check there is enough time between now and revealDeadline

	refundDeadlineMS, err := strconv.ParseUint(refundDeadlineStr, 10, 64)
	if err != nil {
		// xxx
	}
	refundDeadline := bc.FromMillis(refundDeadlineMS)

	// xxx check the time between revealDeadline and refundDeadline isn't too long

	var key [32]byte
	_, err = rand.Read(key[:])
	if err != nil {
		http.Error(w, fmt.Sprintf("choosing cipher key: %s", err), http.StatusInternalServerError)
		return
	}

	rec := &serverRecord{
		clearRoot:      clearRoot,
		key:            key,
		amount:         amount,
		assetID:        assetID,
		revealDeadline: revealDeadline,
		refundDeadline: refundDeadline,
	}

	_, err = rand.Read(rec.transferID[:])
	if err != nil {
		http.Error(w, fmt.Sprintf("choosing transfer ID: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Tedd-Transfer-Id", hex.EncodeToString(rec.transferID[:]))
	w.Header().Set("Content-Type", string(contentType))

	tmpfile, err := ioutil.TempFile("", "teddserve")
	if err != nil {
		http.Error(w, fmt.Sprintf("creating response tempfile: %s", err), http.StatusInternalServerError)
		return
	}
	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)
	defer tmpfile.Close()

	cipherRoot, err := tedd.Serve(tmpfile, f, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("serving data: %s", err), http.StatusInternalServerError)
		return
	}

	err = tmpfile.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("closing response tempfile: %s", err), http.StatusInternalServerError)
		return
	}

	rec.cipherRoot = cipherRoot

	err = s.db.Update(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return errors.Wrap(err, "getting/creating root bucket")
		}
		records, err := root.CreateBucketIfNotExists([]byte("records"))
		if err != nil {
			return errors.Wrap(err, "getting/creating records bucket")
		}
		bu, err := records.CreateBucket(rec.transferID[:])
		if err != nil {
			return errors.Wrapf(err, "creating record bucket %x", rec.transferID[:])
		}
		err = bu.Put([]byte("key"), rec.key[:])
		if err != nil {
			return errors.Wrapf(err, "storing key for record %x", rec.transferID[:])
		}
		err = bu.Put([]byte("clearRoot"), rec.clearRoot)
		if err != nil {
			return errors.Wrapf(err, "storing clearRoot for record %x", rec.transferID[:])
		}
		err = bu.Put([]byte("cipherRoot"), rec.cipherRoot)
		if err != nil {
			return errors.Wrapf(err, "storing cipherRoot for record %x", rec.transferID[:])
		}
		err = bu.Put([]byte("assetID"), rec.assetID)
		if err != nil {
			return errors.Wrapf(err, "storing assetID for record %x", rec.transferID[:])
		}
		var amountBuf [binary.MaxVarintLen64]byte
		m := binary.PutVarint(amountBuf[:], rec.amount)
		err = bu.Put([]byte("amount"), amountBuf[:m])
		if err != nil {
			return errors.Wrapf(err, "storing amount for record %x", rec.transferID[:])
		}
		var revealDeadlineMSBuf [binary.MaxVarintLen64]byte
		m = binary.PutUvarint(revealDeadlineMSBuf[:], bc.Millis(rec.revealDeadline))
		err = bu.Put([]byte("revealDeadlineMS"), revealDeadlineMSBuf[:m])
		if err != nil {
			return errors.Wrapf(err, "storing reveal deadline for record %x", rec.transferID[:])
		}
		var refundDeadlineMSBuf [binary.MaxVarintLen64]byte
		m = binary.PutUvarint(refundDeadlineMSBuf[:], bc.Millis(rec.refundDeadline))
		err = bu.Put([]byte("refundDeadlineMS"), refundDeadlineMSBuf[:m])
		if err != nil {
			return errors.Wrapf(err, "storing refund deadline for record %x", rec.transferID[:])
		}
		return nil
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("storing transfer record: %s", err), http.StatusInternalServerError)
		return
	}

	tmpfile, err = os.Open(tmpfilename)
	if err != nil {
		http.Error(w, fmt.Sprintf("reopening response tempfile: %s", err), http.StatusInternalServerError)
		return
	}
	defer tmpfile.Close()
	_, err = io.Copy(w, tmpfile)
	if err != nil {
		http.Error(w, fmt.Sprintf("writing response: %s", err), http.StatusInternalServerError)
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
		// xxx
	}

	rec, err := s.getRecord(transferID)
	if err != nil {
		// xxx
	}

	paymentProposal, err := hex.DecodeString(paymentProposalStr)
	if err != nil {
		// xxx
	}

	var (
		clearRoot, cipherRoot [32]byte
		assetID               = bc.HashFromBytes(rec.assetID)
	)
	copy(clearRoot[:], rec.clearRoot)
	copy(cipherRoot[:], rec.cipherRoot)

	prog, err := tedd.RevealKey(req.Context(), paymentProposal, s.seller, rec.key, rec.amount, assetID, s.reserver, s.signer, clearRoot, cipherRoot, s.now, rec.revealDeadline, rec.refundDeadline)
	if err != nil {
		// xxx
		return
	}
	vm, err := txvm.Validate(prog, 3, math.MaxInt64)
	if err != nil {
		// xxx
		return
	}
	err = s.submitter(prog, 3, math.MaxInt64-vm.Runlimit())
	if err != nil {
		// xxx
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Runs as a goroutine.
func (s *server) monitorBlockchain(ctx context.Context, url string) {
	client := new(http.Client)
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?height=%d", url, height), nil)
		if err != nil {
			// xxx
		}

		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error getting block at height %d from %s, will retry in ~5 seconds: %s", height, url, err)

			timer := time.NewTimer(5 * time.Second) // xxx add jitter
			select {
			case <-timer.C:
				continue

			case <-ctx.Done():
				// canceled, exit
				timer.Stop()
				log.Print("context canceled, blockchain monitor exiting")
				return
			}
		}

		err = s.processBlock(resp.Body)
		if err != nil {
			// xxx
		}
		if ctx.Err() != nil {
			// canceled, exit
			log.Print("context canceled, blockchain monitor exiting")
			return
		}
	}
}

func (s *server) processBlock(r io.ReadCloser) error {
	defer r.Close()

	bits, err := ioutil.ReadAll(r)
	if err != nil {
		// xxx
	}
	b := new(bc.Block)
	err = b.FromBytes(bits)
	if err != nil {
		// xxx
	}
	// xxx set server's blockchain timestamp
	return s.reserver.processBlock(b)
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
		copy(rec.key[:], bu.Get([]byte("key")))
		rec.clearRoot = bu.Get([]byte("clearRoot"))
		rec.cipherRoot = bu.Get([]byte("cipherRoot"))
		rec.assetID = bu.Get([]byte("assetID"))

		var n int
		rec.amount, n = binary.Varint(bu.Get([]byte("amount")))
		if n < 1 {
			return fmt.Errorf("cannot parse amount in record %x", transferID)
		}
		revealDeadlineMS, n := binary.Uvarint(bu.Get([]byte("revealDeadlineMS")))
		if n < 1 {
			return fmt.Errorf("cannot parse reveal deadline in record %x", transferID)
		}
		rec.revealDeadline = bc.FromMillis(revealDeadlineMS)
		refundDeadlineMS, n := binary.Uvarint(bu.Get([]byte("refundDeadlineMS")))
		if n < 1 {
			return fmt.Errorf("cannot parse refund deadline in record %x", transferID)
		}
		rec.refundDeadline = bc.FromMillis(refundDeadlineMS)
		return nil
	})
	return &rec, err
}
