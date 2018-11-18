package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bobg/merkle"
	"github.com/bobg/tredd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txvm"
	"github.com/coreos/bbolt"
)

func get(args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		clearRootHex         = fs.String("hash", "", "clear-chunk Merkle root hash of requested file")
		amount               = fs.Int64("amount", 0, "amount of proposed payment")
		assetIDHex           = fs.String("asset", "", "asset ID of proposed payment")
		revealDeadlineDurStr = fs.String("reveal", "", "time until reveal deadline, in time.ParseDuration format")
		refundDeadlineDurStr = fs.String("refund", "", "time from reveal deadline until refund deadline")
		dbFile               = fs.String("db", "", "file containing client-state db")
		prvFile              = fs.String("prv", "", "file containing client private key")
		serverURL            = fs.String("server", "", "base URL of tredd server")
		bcURL                = fs.String("bcurl", "", "base URL of blockchain server")
		dir                  = fs.String("dir", "", "root dir for file transfers")
	)

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	var (
		requestURL        = *serverURL + "/request"
		proposePaymentURL = *serverURL + "/propose-payment"
	)

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
	buyer := prv.Public().(ed25519.PublicKey)

	var clearRoot [32]byte
	_, err = hex.Decode(clearRoot[:], []byte(*clearRootHex))
	if err != nil {
		log.Fatal(err)
	}

	assetIDBytes, err := hex.DecodeString(*assetIDHex)
	if err != nil {
		log.Fatal(err)
	}
	assetID := bc.HashFromBytes(assetIDBytes)

	revealDeadlineDur, err := time.ParseDuration(*revealDeadlineDurStr)
	if err != nil {
		log.Fatal(err)
	}
	revealDeadline := time.Now().Add(revealDeadlineDur)

	refundDeadlineDur, err := time.ParseDuration(*refundDeadlineDurStr)
	if err != nil {
		log.Fatal(err)
	}
	refundDeadline := revealDeadline.Add(refundDeadlineDur)

	db, err := bbolt.Open(*dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("launching blockchain observer")
	o := newObserver(db, buyer, *bcURL+"/get")

	q := newQuiescenceWaiter()
	o.setcb(func(*bc.Tx) { q.ping() })
	go o.run(ctx)

	log.Print("waiting to catch up...")
	q.wait()
	log.Print("...caught up")

	vals := url.Values{}
	vals.Add("clearroot", *clearRootHex)
	vals.Add("amount", strconv.FormatInt(*amount, 10))
	vals.Add("assetid", *assetIDHex)
	vals.Add("revealdeadline", strconv.FormatInt(int64(bc.Millis(revealDeadline)), 10)) // xxx range check
	vals.Add("refunddeadline", strconv.FormatInt(int64(bc.Millis(refundDeadline)), 10)) // xxx range check

	log.Print("requesting content")
	resp, err := http.PostForm(requestURL, vals)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatalf("status code %d from initial HTTP request", resp.StatusCode)
	}

	var (
		transferID       = resp.Header.Get("X-Tredd-Transfer-Id")
		clearHashesFile  = path.Join(*dir, fmt.Sprintf("hashes-%s", transferID))
		cipherChunksFile = path.Join(*dir, fmt.Sprintf("chunks-%s", transferID))
	)

	clearHashes, err := newFileChunkStore(clearHashesFile, 32)
	if err != nil {
		log.Fatalf("creating hash chunk store: %s", err)
	}
	defer os.Remove(clearHashesFile) // TODO: keep this around if needed to recover from errors

	cipherChunks, err := newFileChunkStore(cipherChunksFile, tredd.ChunkSize)
	if err != nil {
		log.Fatalf("creating cipher chunk store: %s", err)
	}
	defer os.Remove(cipherChunksFile) // TODO: keep this around if needed to recover from errors

	log.Print("storing cipher chunks and checking clear hashes")
	cipherRoot, err := tredd.Get(resp.Body, clearRoot, clearHashes, cipherChunks)
	if err != nil {
		log.Fatal(err)
	}

	signer := func(msg []byte) ([]byte, error) {
		return ed25519.Sign(prv, msg), nil
	}

	var cipherRootBuf [32]byte
	copy(cipherRootBuf[:], cipherRoot)

	now := time.Now()

	prog, err := tredd.ProposePayment(ctx, buyer, *amount, assetID, clearRoot, cipherRootBuf, now, revealDeadline, refundDeadline, o.r, signer)
	if err != nil {
		log.Fatal(err)
	}

	parsed := tredd.ParseLog(prog)
	if parsed == nil {
		log.Fatal("cannot parse log of proposed payment transaction")
	}
	anchor1 := parsed.Anchor1

	submit := submitter(*bcURL + "/submit")

	o.setcb(func(tx *bc.Tx) {
		parsed := tredd.ParseLog(tx.Program)
		if parsed == nil {
			return
		}
		if !bytes.Equal(parsed.Anchor1, anchor1) {
			return
		}

		defer cancel()

		log.Printf("payment proposal accepted, key is %x; now decrypting", parsed.Key)

		// Payment has been accepted.
		var key [32]byte
		copy(key[:], parsed.Key)

		outFileName := path.Join(*dir, hex.EncodeToString(clearRoot[:]))
		out, err := os.Create(outFileName)
		if err != nil {
			log.Fatalf("creating %s: %s", outFileName, err) // TODO: more graceful/recoverable handling
		}
		defer out.Close()

		err = tredd.Decrypt(out, clearHashes, cipherChunks, key)
		if bchErr, ok := err.(tredd.BadClearHashError); ok {
			log.Printf("decryption failed on chunk %d; now claiming refund", bchErr.Index)

			redeem := &tredd.Redeem{
				RefundDeadline: refundDeadline,
				Buyer:          buyer,
				Seller:         parsed.Seller,
				Amount:         2 * *amount,
				AssetID:        assetID,
				ClearRoot:      clearRoot,
				Key:            key,
			}
			copy(redeem.CipherRoot[:], cipherRoot)
			copy(redeem.Anchor2[:], parsed.Anchor2)

			var (
				refHash        [32 + binary.MaxVarintLen64]byte
				refCipherChunk [tredd.ChunkSize + binary.MaxVarintLen64]byte
			)
			m := binary.PutUvarint(refHash[:], bchErr.Index)
			binary.PutUvarint(refCipherChunk[:], bchErr.Index)

			g, err := clearHashes.Get(bchErr.Index)
			if err != nil {
				log.Fatalf("getting hash %d from %s: %s", bchErr.Index, clearHashes.filename, err)
			}
			copy(refHash[m:], g)

			g, err = cipherChunks.Get(bchErr.Index)
			if err != nil {
				log.Fatalf("getting cipher chunk %d from %s: %s", bchErr.Index, cipherChunks.filename, err)
			}
			copy(refCipherChunk[m:], g)

			var (
				clearTree  = merkle.NewProofTree(sha256.New(), refHash[:m+32])
				cipherTree = merkle.NewProofTree(sha256.New(), refCipherChunk[:m+len(g)])
				hasher     = sha256.New()
			)
			nchunks, err := cipherChunks.Len()
			if err != nil {
				log.Fatalf("getting length of cipher chunk store %s: %s", cipherChunks.filename, err)
			}
			for index := uint64(0); index < uint64(nchunks); index++ {
				var chunk [tredd.ChunkSize + binary.MaxVarintLen64]byte
				m := binary.PutUvarint(chunk[:], index)
				ci, err := cipherChunks.Get(index)
				if err != nil {
					log.Fatalf("getting cipher chunk %d from %s: %s", bchErr.Index, cipherChunks.filename, err)
				}
				copy(chunk[m:], ci)
				n := len(ci)

				var h [32 + binary.MaxVarintLen64]byte
				binary.PutUvarint(h[:], index)
				merkle.LeafHash(hasher, h[:m], chunk[:m+n])

				clearTree.Add(h[:m+32])
				cipherTree.Add(chunk[:m+n])
			}

			var (
				clearProof  = clearTree.Proof()
				cipherProof = cipherTree.Proof()
			)

			prog, err := tredd.ClaimRefund(redeem, int64(bchErr.Index), refCipherChunk[m:m+len(g)], refHash[m:m+32], cipherProof, clearProof) // xxx range check
			if err != nil {
				log.Fatalf("constructing refund-claiming transaction: %s", err)
			}

			vm, err := txvm.Validate(prog, 3, math.MaxInt64)
			if err != nil {
				log.Fatalf("calculating runlimit for refund-claiming transaction: %s", err)
			}

			err = submit(prog, 3, math.MaxInt64-vm.Runlimit())
			if err != nil {
				// TODO: retry
				log.Fatalf("submitting refund-claiming transaction: %s", err)
			}
			return
		}
		if err != nil {
			log.Fatalf("decrypting content: %s", err)
		}
		log.Print("complete")
	})
	o.enqueue(revealDeadline, func() {
		log.Print("reveal deadline has arrived, transfer invalid")
		cancel()
	})

	log.Print("proposing payment")
	req, err := http.NewRequest("POST", proposePaymentURL, bytes.NewReader(prog))
	if err != nil {
		log.Fatalf("constructing payment proposal: %s", err)
	}
	req = req.WithContext(ctx)

	req.Header.Set("X-Tredd-Transfer-Id", transferID)

	var client http.Client
	resp, err = client.Do(req) // from this point, funds are committed - perhaps even in case of error
	if err != nil {
		log.Printf("sending payment proposal: %s", err)
		log.Print("WARNING: funds may be committed; awaiting outcome")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("sending payment proposal: unexpected status %d", resp.StatusCode)
		log.Print("WARNING: funds may be committed; awaiting outcome")
	}

	log.Print("awaiting key or reveal deadline")
	<-ctx.Done()
}

type quiescenceWaiter struct {
	mu sync.Mutex
	c  *sync.Cond
	t  time.Time
}

func newQuiescenceWaiter() *quiescenceWaiter {
	q := new(quiescenceWaiter)
	q.c = sync.NewCond(&q.mu)
	return q
}

func (q *quiescenceWaiter) ping() {
	q.mu.Lock()
	q.t = time.Now()
	q.c.Broadcast()
	q.mu.Unlock()
}

// wait returns when q has not been pinged for one second.
func (q *quiescenceWaiter) wait() {
	ch := make(chan struct{})
	var done int32
	go func() {
		q.mu.Lock()
		defer q.mu.Unlock()
		for atomic.LoadInt32(&done) == 0 {
			ch <- struct{}{}
			q.c.Wait()
		}
		close(ch)
	}()
	for {
		t := time.NewTimer(time.Second)
		select {
		case <-t.C:
			atomic.StoreInt32(&done, 1)
			return
		case <-ch:
			t.Stop()
		}
	}
}
