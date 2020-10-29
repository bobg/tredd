package main

import (
	"context"
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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/bobg/tredd/contract"

	"github.com/bobg/tredd"
)

func get(args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		clearRootHex      = fs.String("hash", "", "clear-chunk Merkle root hash of requested file")
		tokenTypeStr      = fs.String("token", "", "token type (ERC20 hex address) of proposed payment, or omit for ETH")
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
		proposeURL     = *serverURL + "/propose-payment"
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

	var tokenType common.Address
	if *tokenTypeStr != "" {
		tokenType = common.HexToAddress(*tokenTypeStr)
	}

	vals := url.Values{}
	vals.Add("buyer", buyer.From.Hex())
	vals.Add("clearroot", *clearRootHex)
	vals.Add("amount", amount.String())
	vals.Add("collateral", collateral.String())
	vals.Add("revealdeadline", strconv.FormatInt(revealDeadline.Unix(), 10))
	vals.Add("refunddeadline", strconv.FormatInt(refundDeadline.Unix(), 10)) // TODO: range check
	if tokenType != (common.Address{}) {
		vals.Add("token", tokenType.Hex())
	}

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

	contractAddr, con, _, err := tredd.ProposePayment(ctx, client, buyer, seller, tokenType, amount, collateral, clearRoot, cipherRootBuf, revealDeadline, refundDeadline)
	if err != nil {
		log.Fatal(err)
	}

	vals = url.Values{}
	vals.Add("transferid", transferID)
	vals.Add("contractaddr", contractAddr.Hex())
	resp, err = http.PostForm(proposeURL, vals)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	evChan := make(chan *contract.TreddEvDecryptionKey)
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

			refClearHash, refCipherChunk, clearProof, cipherProof, err := tredd.PrepareForRefund(bchErr.Index, clearHashes, cipherChunks)
			if err != nil {
				log.Fatalf("preparing for refund: %s", err)
			}

			receipt, err := tredd.ClaimRefund(ctx, client, buyer, con, bchErr.Index, refCipherChunk, refClearHash, cipherProof, clearProof)
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
