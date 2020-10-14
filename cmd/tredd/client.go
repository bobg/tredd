package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
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
		tokenTypeStr      = fs.String("token", "", "token type (ERC20 hex address) of proposed payment")
		amountStr         = fs.String("amount", "1", "amount of proposed payment")
		collateralStr     = fs.String("collateral", "1", "amount of proposed collateral")
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

	var (
		amount     = new(big.Int)
		collateral = new(big.Int)
	)
	_, ok := amount.SetString(*amountStr, 10)
	if !ok {
		log.Fatalf(`Error parsing amount string "%s"`, *amountStr)
	}
	_, ok = collateral.SetString(*collateralStr, 10)
	if !ok {
		log.Fatalf(`Error parsing collateralStr string "%s"`, *collateralStr)
	}

	client, err := ethclient.Dial(*ethURL)
	if err != nil {
		log.Fatal(err)
	}

	tokenType := common.HexToAddress(*tokenTypeStr)

	vals := url.Values{}
	vals.Add("clearroot", *clearRootHex)
	vals.Add("amount", amount.String())
	vals.Add("collateral", collateral.String())
	vals.Add("token", tokenType.Hex())
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

	var seller common.Address
	_, err = hex.Decode(seller[:], []byte(*sellerHex))
	if err != nil {
		log.Fatal(err)
	}

	con, err := tredd.ProposePayment(ctx, client, buyer, seller, tokenType, amount, collateral, clearRoot, cipherRootBuf, revealDeadline, refundDeadline)
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
		receipt, err := tredd.Cancel(ctx, client, buyer, con)
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

			refClearHash, err := clearHashes.Get(bchErr.Index)
			if err != nil {
				log.Fatalf("Error getting clear hash %d: %s", bchErr.Index, err)
			}
			var refClearHashBuf [32]byte
			copy(refClearHashBuf[:], refClearHash)
			prefixedRefClearHash, err := tredd.PrefixHash(uint64(bchErr.Index), refClearHashBuf)
			if err != nil {
				log.Fatalf("Error prefixing clear hash %d: %s", bchErr.Index, err)
			}

			refCipherChunk, err := cipherChunks.Get(bchErr.Index)
			if err != nil {
				log.Fatalf("Error getting cipher chunk %d: %s", bchErr.Index, err)
			}
			prefixedRefCipherChunk, err := tredd.PrefixChunk(uint64(bchErr.Index), refCipherChunk)
			if err != nil {
				log.Fatalf("Error prefixing cipher chunk %d: %s", bchErr.Index, err)
			}

			var (
				clearTree  = merkle.NewProofTree(sha256.New(), prefixedRefClearHash)
				cipherTree = merkle.NewProofTree(sha256.New(), prefixedRefCipherChunk)
			)

			nchunks, err := cipherChunks.Len()
			if err != nil {
				log.Fatalf("Error getting size of cipher-chunk store: %s", err)
			}

			for index := int64(0); index < nchunks; index++ {
				clearHash, err := clearHashes.Get(index)
				if err != nil {
					log.Fatalf("Error getting clear hash %d: %s", index, err)
				}
				var clearHashBuf [32]byte
				copy(clearHashBuf[:], clearHash)
				prefixedClearHash, err := tredd.PrefixHash(uint64(index), clearHashBuf)
				if err != nil {
					log.Fatalf("Error prefixing clear hash %d: %s", index, err)
				}

				cipherChunk, err := cipherChunks.Get(index)
				if err != nil {
					log.Fatalf("Error getting cipher chunk %d: %s", index, err)
				}
				prefixedCipherChunk, err := tredd.PrefixChunk(uint64(index), cipherChunk)
				if err != nil {
					log.Fatalf("Error prefixing cipher chunk %d: %s", index, err)
				}

				clearTree.Add(prefixedClearHash)
				cipherTree.Add(prefixedCipherChunk)
			}

			var (
				clearProof  = clearTree.Proof()
				cipherProof = cipherTree.Proof()
			)

			receipt, err := tredd.ClaimRefund(ctx, client, buyer, con, bchErr.Index, refCipherChunk, refClearHashBuf, cipherProof, clearProof)
			if err != nil {
				log.Fatalf("Error constructing refund-claiming transaction: %s", err)
			}

			log.Printf("Refund claimed in transaction %x", receipt.TxHash[:])
			return

		} else if err != nil {
			log.Fatalf("Error decrypting content: %s", err)
		}
		log.Printf("Complete, decrypted content is in %s", outFileName)

	case err := <-subErrChan:
		log.Fatalf("Error waiting for decryption-key event: %s", err)
	}
}
