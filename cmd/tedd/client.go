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
	"time"

	"github.com/bobg/merkle"
	"github.com/bobg/tedd"
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
		clearRootHex         = flag.String("hash", "", "clear-chunk Merkle root hash of requested file")
		amount               = flag.Int64("amount", 0, "amount of proposed payment")
		assetIDHex           = flag.String("asset", "", "asset ID of proposed payment")
		revealDeadlineDurStr = flag.String("reveal", "", "time until reveal deadline, in time.ParseDuration format")
		refundDeadlineDurStr = flag.String("refund", "", "time from reveal deadline until refund deadline")
		dbFile               = fs.String("db", "", "file containing client-state db")
		prvFile              = fs.String("prv", "", "file containing client private key")
		serverURL            = fs.String("server", "", "base URL of tedd server")
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

	o := newObserver(db, buyer, *bcURL+"/get")
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

	var (
		transferID       = resp.Header.Get("X-Tedd-Transfer-Id")
		clearHashesFile  = path.Join(*dir, fmt.Sprintf("hashes-%s", transferID))
		cipherChunksFile = path.Join(*dir, fmt.Sprintf("chunks-%s", transferID))
	)

	clearHashes := &fileChunkStore{
		filename:  clearHashesFile,
		chunksize: 32,
	}
	defer os.Remove(clearHashesFile)

	cipherChunks := &fileChunkStore{
		filename:  cipherChunksFile,
		chunksize: tedd.ChunkSize,
	}
	defer os.Remove(cipherChunksFile)

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

	submit := submitter(*bcURL + "/submit")

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

		out, err := os.Create(path.Join(*dir, hex.EncodeToString(clearRoot[:])))
		if err != nil {
			// xxx
		}
		defer out.Close()

		err = tedd.Decrypt(out, clearHashes, cipherChunks, key)
		if bchErr, ok := err.(tedd.BadClearHashError); ok {
			log.Print("payment accepted but decryption failed, now claiming refund")

			redeem := &tedd.Redeem{
				RefundDeadline: refundDeadline,
				Buyer:          buyer,
				Seller:         parsed.Seller,
				Amount:         *amount, // xxx right?
				AssetID:        assetID,
				ClearRoot:      clearRoot,
				Key:            key,
			}
			copy(redeem.CipherRoot[:], cipherRoot)
			copy(redeem.Anchor[:], parsed.Anchor2) // xxx right?

			var (
				refHash        [32 + binary.MaxVarintLen64]byte
				refCipherChunk [tedd.ChunkSize + binary.MaxVarintLen64]byte
			)
			m := binary.PutUvarint(refHash[:], bchErr.Index)
			binary.PutUvarint(refCipherChunk[:], bchErr.Index)

			g, err := clearHashes.Get(bchErr.Index)
			if err != nil {
				// xxx
			}
			copy(refHash[m:], g)

			g, err = cipherChunks.Get(bchErr.Index)
			if err != nil {
				// xxx
			}
			copy(refCipherChunk[m:], g)

			var (
				clearTree  = merkle.NewProofTree(sha256.New(), refHash[:m+32])
				cipherTree = merkle.NewProofTree(sha256.New(), refCipherChunk[:m+len(g)])
				hasher     = sha256.New()
			)
			nchunks, err := cipherChunks.Len()
			if err != nil {
				// xxx
			}
			for index := uint64(0); index < uint64(nchunks); index++ {
				var chunk [tedd.ChunkSize + binary.MaxVarintLen64]byte
				m := binary.PutUvarint(chunk[:], index)
				ci, err := cipherChunks.Get(index)
				if err != nil {
					// xxx
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

			prog, err := tedd.ClaimRefund(redeem, int64(bchErr.Index), refCipherChunk[m:m+len(g)], refHash[m:m+32], cipherProof, clearProof) // xxx range check
			if err != nil {
				// xxx
			}

			vm, err := txvm.Validate(prog, 3, math.MaxInt64)
			if err != nil {
				// xxx
			}

			err = submit(prog, 3, math.MaxInt64-vm.Runlimit())
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

	req, err := http.NewRequest("POST", proposePaymentURL, bytes.NewReader(prog))
	if err != nil {
		// xxx
	}
	req = req.WithContext(ctx)

	req.Header.Set("X-Tedd-Transfer-Id", transferID)

	var client http.Client
	resp, err = client.Do(req)
	if err != nil {
		// xxx
	}
	defer resp.Body.Close()

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
