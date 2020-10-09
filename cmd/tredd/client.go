package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/bobg/tredd"
)

func get(args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		clearRootHex      = fs.String("hash", "", "clear-chunk Merkle root hash of requested file")
		amount            = fs.Int64("amount", 0, "amount of proposed payment")
		tokenType         = fs.String("token", "", "asset ID of proposed payment")
		revealDeadlineDur = fs.Duration("reveal", 15*time.Minute, "time until reveal deadline, in time.ParseDuration format")
		refundDeadlineDur = fs.Duration("refund", 30*time.Minute, "time from reveal deadline until refund deadline, in time.ParseDuration format")
		serverURL         = fs.String("server", "", "base URL of tredd server")
		ethURL            = fs.String("ethurl", "", "base URL of Ethereum server")
		dir               = fs.String("dir", "", "root dir for file transfers")
		sellerHex         = fs.String("seller", "", "seller address (hex)")
	)

	keyfile, passphrase := addKeyfilePassphrase(fs)

	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	var (
		requestURL     = *serverURL + "/request"
		revealDeadline = time.Now().Add(*revealDeadlineDur)
		refundDeadline = revealDeadline.Add(*refundDeadlineDur)
	)

	var clearRoot [32]byte
	_, err = hex.Decode(clearRoot[:], []byte(*clearRootHex))
	if err != nil {
		log.Fatal(err)
	}

	buyer, err := handleKeyfilePassphrase(*keyfile, *passphrase)
	if err != nil {
		log.Fatal(err)
	}

	vals := url.Values{}
	vals.Add("clearroot", *clearRootHex)
	vals.Add("amount", strconv.FormatInt(*amount, 10))
	vals.Add("token", *tokenType)
	vals.Add("revealdeadline", strconv.FormatInt(int64(Millis(revealDeadline)), 10)) // TODO: range check
	vals.Add("refunddeadline", strconv.FormatInt(int64(Millis(refundDeadline)), 10)) // TODO: range check

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

	var cipherRootBuf [32]byte
	copy(cipherRootBuf[:], cipherRoot)

	log.Print("proposing payment")

	client, err := ethclient.Dial(*ethURL)
	if err != nil {
		log.Fatal(err)
	}

	var seller common.Address
	_, err = hex.Decode(seller[:], []byte(*sellerHex))
	if err != nil {
		log.Fatal(err)
	}

	receipt, err := tredd.ProposePayment(ctx, client, buyer, seller, *amount, *tokenType, clearRoot, cipherRootBuf, revealDeadline, refundDeadline)
	if err != nil {
		log.Fatal(err)
	}
	contractAddr := receipt.ContractAddress
	con, err := tredd.NewTredd(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	evChan := make(chan *tredd.TreddEvDecryptionKey)
	sub, err := con.WatchEvDecryptionKey(&bind.WatchOpts{Context: ctx}, evChan)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()
	subErrChan := sub.Err()

	revealTimer := time.NewTimer(time.Until(revealDeadline))
	defer revealTimer.Stop()

	// Wait for the reveal deadline to pass,
	// in which case we reclaim payment from the contract,
	// or for the reveal-key event,
	// in which case we decrypt and validate the content.
	select {
	case <-ctx.Done():
		log.Print("context canceled, exiting")
		return

	case <-revealTimer.C:
		// Deadline passed with no decryption key revealed.
		receipt, err := con.ReclaimPayment(ctx, client, buyer)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Reclaimed payment in transaction %x", receipt.TxHash[:])
		return

	case ev := <-evChan:
		// Decryption key revealed.
		outFileName := path.Join(*dir, hex.EncodeToString(clearRoot[:]))
		out, err := os.Create(outFileName)
		if err != nil {
			log.Fatalf("creating %s: %s", outFileName, err) // TODO: more graceful/recoverable handling
		}
		defer out.Close()

		var bchErr tredd.BadClearHashError
		err = tredd.Decrypt(out, clearHashes, cipherChunks, ev.DecryptionKey)
		if errors.As(err, &bchErr) {
			// Validation failed, claim a refund.

			log.Printf("decryption failed on chunk %d; now claiming refund", bchErr.Index)

			var (
				refHash        [32 + binary.MaxVarintLen64]byte
				refCipherChunk [tredd.ChunkSize + binary.MaxVarintLen64]byte
			)
			m := binary.PutUvarint(refHash[:], uint64(bchErr.Index))
			binary.PutUvarint(refCipherChunk[:], uint64(bchErr.Index))

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
			for index := int64(0); index < int64(nchunks); index++ {
				var chunk [tredd.ChunkSize + binary.MaxVarintLen64]byte
				m := binary.PutUvarint(chunk[:], uint64(index))
				ci, err := cipherChunks.Get(index)
				if err != nil {
					log.Fatalf("getting cipher chunk %d from %s: %s", bchErr.Index, cipherChunks.filename, err)
				}
				copy(chunk[m:], ci)
				n := len(ci)

				var h [32 + binary.MaxVarintLen64]byte
				binary.PutUvarint(h[:], uint64(index))
				merkle.LeafHash(hasher, h[:m], chunk[:m+n])

				clearTree.Add(h[:m+32])
				cipherTree.Add(chunk[:m+n])
			}

			var (
				clearProof  = clearTree.Proof()
				cipherProof = cipherTree.Proof()
				clearHash   [32]byte
			)

			copy(clearHash[:], refHash[m:m+32])

			receipt, err := tredd.ClaimRefund(ctx, client, buyer, contractAddr, bchErr.Index, refCipherChunk[m:m+len(g)], clearHash, cipherProof, clearProof) // TODO: range check
			if err != nil {
				log.Fatalf("Error constructing refund-claiming transaction: %s", err)
			}

			log.Printf("Refund claimed in transaction %x", receipt.TxHash[:])
			return

		} else if err != nil {
			log.Fatalf("Error decryption content: %s", err)
		}
		log.Printf("Complete, decrypted content is in %s", outFileName)

	case err := <-subErrChan:
		log.Fatalf("Error waiting for decryption-key event: %s", err)
	}
}
