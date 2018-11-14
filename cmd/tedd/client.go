package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bobg/tedd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/protocol/bc"
	"github.com/coreos/bbolt"
)

func get(args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		clearRootHex         = flag.String("hash", "", "clear-chunk Merkle root hash of requested file")
		amount               = flag.Int64("amount", 0, "amount of proposed payment")
		assetIDHex           = flag.String("asset", "", "asset ID of proposed payment")
		revealDeadlineDurStr = flag.String("reveal", "", "time until reveal deadline, in time.ParseDuration format")
		refundDeadlineDurStr = flag.String("refund", "", "time from reveal deadline until refund deadline")
		dbFile               = fs.String("db", "", "file containing client-state db")
		prvFile              = fs.String("prv", "", "file containing client private key")
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

	o := newObserver(db, buyer)
	go o.run(ctx)

	var vals url.Values
	vals.Add("clearroot", *clearRootHex)
	vals.Add("amount", strconv.FormatInt(*amount, 10))
	vals.Add("assetid", *assetIDHex)
	vals.Add("revealdeadline", strconv.FormatInt(int64(bc.Millis(revealDeadline)), 10)) // xxx range check
	vals.Add("refunddeadline", strconv.FormatInt(int64(bc.Millis(refundDeadline)), 10)) // xxx range check

	resp, err := http.PostForm(requestURL, vals)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatalf("status code %d from initial HTTP request", resp.StatusCode)
	}

	cipherRoot, err := tedd.Get(&bytereader{r: resp.Body}, clearRoot, clearHashes, cipherChunks)
	if err != nil {
		log.Fatal(err)
	}

	signer := func(msg []byte) ([]byte, error) {
		return ed25519.Sign(prv, msg), nil
	}

	var cipherRootBuf [32]byte
	copy(cipherRootBuf[:], cipherRoot)

	now, err := o.now()
	if err != nil {
		log.Fatal(err)
	}

	prog, err := tedd.ProposePayment(ctx, buyer, *amount, assetID, clearRoot, cipherRootBuf, now, revealDeadline, refundDeadline, o.r, signer)
	if err != nil {
		log.Fatal(err)
	}

	parsed := tedd.ParseLog(prog)
	if parsed == nil {
		log.Fatal("cannot parse log of proposed payment transaction")
	}
	anchor1 := parsed.Anchor1

	o.setcb(func(tx *bc.Tx) {
		defer cancel()

		parsed := tedd.ParseLog(tx.Program)
		if parsed == nil {
			return
		}
		if !bytes.Equal(parsed.Anchor1, anchor1) {
			return
		}
		// Payment has been accepted.
		var key [32]byte
		copy(key[:], parsed.Key)
		err := tedd.Decrypt(out, clearHashes, cipherChunks, key)
		if err, ok := err.(tedd.BadClearHashError); ok {
			log.Print("payment accepted but decryption failed, now claiming refund")

			redeem := &tedd.Redeem{
				// xxx
			}
			prog, err := tedd.ClaimRefund(redeem, err.Index, cipherChunk, clearHash, cipherProof, clearProof)
			if err != nil {
				// xxx
			}
			err = submitter(xxx)
			if err != nil {
				// xxx
			}
			return
		}
		if err != nil {
			// xxx
		}
		log.Print("complete")
	})
	o.enqueue(revealDeadline, func() {
		log.Print("reveal deadline has arrived, transfer invalid")
		// xxx remove encrypted file
		cancel()
	})

	resp, err = http.Post(proposePaymentURL, "application/octet-stream", bytes.NewReader(prog))
	if err != nil {
		// xxx
	}

	if resp.StatusCode != http.StatusNoContent {
		// xxx
	}

	<-ctx.Done()
}

type bytereader struct {
	r io.Reader
}

func (b *bytereader) ReadByte() (byte, error) {
	var buf [1]byte
	n, err := b.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, io.ErrUnexpectedEOF
	}
	return buf[0], nil
}

func (b *bytereader) Read(buf []byte) (int, error) {
	return b.r.Read(buf)
}
