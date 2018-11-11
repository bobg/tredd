package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/bobg/tedd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txvm"
)

func serve(args []string) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		prvFile = fs.String("prv", "", "file containing server private key")
	)

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

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

	s := &server{
		dir:    dir,
		seller: prv.Public().(ed25519.PublicKey),
	}
	s.signer = func(msg []byte) ([]byte, error) {
		return ed25519.Sign(prv, msg), nil
	}

	// xxx queue claim-payment calls for existing records

	go s.monitorBlockchain(bcGetURL)

	http.HandleFunc("/", s.serve)
	http.HandleFunc("/propose-payment", s.acceptPayment)
	http.ListenAndServe(addr, nil)
}

type server struct {
	dir       string
	seller    ed25519.PublicKey
	reserver  *reserver // must satisfy tedd.Reserver
	signer    tedd.Signer
	submitter func(prog []byte, version, runlimit int64) error
}

type serverRecord struct {
	clearRoot, cipherRoot [32]byte
	key                   [32]byte
	amount                int64
	assetID               bc.Hash
	refundDeadline        time.Time
}

func (s *server) serve(w http.ResponseWriter, req *http.Request) {
	// xxx parse request
	// xxx check revealDeadline is far enough in the future, and refundDeadline is soon enough after that
	// xxx check amount/assetID is acceptable and that there's enough for collateral

	var key [32]byte
	_, err := rand.Read(key[:])
	if err != nil {
		http.Error(w, fmt.Sprintf("choosing cipher key: %s", err), http.StatusInternalServerError)
		return
	}

	rec := &serverRecord{
		clearRoot:      clearRoot,
		key:            key,
		amount:         amount,
		assetID:        assetID,
		refundDeadline: refundDeadline,
	}

	// xxx set header fields

	f, err := os.Open(path.Join(s.dir, filename))
	if os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("opening %s: %s", filename, err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	cipherRoot, err := tedd.Serve(w, f, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("serving data: %s", err), http.StatusInternalServerError)
		return
	}

	rec.cipherRoot = cipherRoot

	// xxx store record
}

func (s *server) revealKey(w http.ResponseWriter, req *http.Request) {
	// xxx parse request
	prog, err := tedd.RevealKey(req.Context(), paymentProposal, s.seller, key, amount, assetID, s.reserver, s.signer, clearRoot, revealDeadline, refundDeadline)
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
}

// Runs as a goroutine.
func (s *server) monitorBlockchain(ctx context.Context, url string) {
	client := new(http.Client)
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s?height=%d", url, height))
		if err != nil {
			log.Printf("error getting block at height %d from %s, will retry in ~5 seconds: %s", height, url, err)

			timer := time.Timer(5 * time.Second) // xxx add jitter
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

		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			// xxx
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

	bits, err := ioutil.ReadAll(resp.Body)
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
