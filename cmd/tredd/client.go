package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bobg/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bobg/tredd"
	"github.com/bobg/tredd/contract"
)

func get(
	ctx context.Context,
	clearRootHex string,
	tokenTypeStr string,
	amountStr string,
	collateralStr string,
	revealDeadlineDur time.Duration,
	refundDeadlineDur time.Duration,
	serverURL string,
	ethURL string,
	dir string,
	sellerHex string,
	keyfile string,
	passphrase string,
	_ []string,
) error {
	var (
		requestURL     = serverURL + "/request"
		proposeURL     = serverURL + "/propose-payment"
		revealDeadline = time.Now().Add(revealDeadlineDur)
		refundDeadline = revealDeadline.Add(refundDeadlineDur)
	)

	var clearRoot [32]byte
	if _, err := hex.Decode(clearRoot[:], []byte(clearRootHex)); err != nil {
		return errors.Wrap(err, "decoding clear root hex")
	}

	buyer, err := handleKeyfilePassphrase(keyfile, passphrase)
	if err != nil {
		return errors.Wrap(err, "handling keyfile and passphrase")
	}

	var (
		amount     = new(big.Int)
		collateral = new(big.Int)
	)
	if _, ok := amount.SetString(amountStr, 10); !ok {
		return fmt.Errorf("error parsing amount string %q", amountStr)
	}
	if _, ok := collateral.SetString(collateralStr, 10); !ok {
		return fmt.Errorf("error parsing collateralStr string %q", collateralStr)
	}

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return errors.Wrapf(err, "dialing Ethereum service at %s", ethURL)
	}

	var tokenType common.Address
	if tokenTypeStr != "" {
		tokenType = common.HexToAddress(tokenTypeStr)
	}

	vals := url.Values{}
	vals.Add("buyer", buyer.From.Hex())
	vals.Add("clearroot", clearRootHex)
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
		clearHashesFile  = filepath.Join(dir, fmt.Sprintf("hashes-%s", transferID))
		cipherChunksFile = filepath.Join(dir, fmt.Sprintf("chunks-%s", transferID))
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
	_, err = hex.Decode(seller[:], []byte(sellerHex))
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
		return context.Cause(ctx)

	case <-revealTimer.C:
		receipt, err := tredd.Cancel(ctx, client, buyer, con)
		if err != nil {
			return errors.Wrap(err, "canceling contract after reveal deadline")
		}
		log.Printf("Reclaimed payment in transaction %x", receipt.TxHash[:])
		return nil

	case ev := <-evChan:
		// Decryption key revealed.
		outFileName := filepath.Join(dir, hex.EncodeToString(clearRoot[:]))
		out, err := os.Create(outFileName)
		if err != nil {
			return errors.Wrapf(err, "creating output file %s", outFileName) // TODO: more graceful/recoverable handling
		}
		defer out.Close()

		var bchErr tredd.BadClearHashError
		err = tredd.Decrypt(out, clearHashes, cipherChunks, ev.DecryptionKey)
		if errors.As(err, &bchErr) {
			// Validation failed, claim a refund.

			log.Printf("decryption failed on chunk %d; now claiming refund", bchErr.Index)

			refClearHash, refCipherChunk, clearProof, cipherProof, err := tredd.PrepareForRefund(bchErr.Index, clearHashes, cipherChunks)
			if err != nil {
				return errors.Wrap(err, "preparing for refund")
			}

			receipt, err := tredd.ClaimRefund(ctx, client, buyer, con, bchErr.Index, refCipherChunk, refClearHash, cipherProof, clearProof)
			if err != nil {
				return errors.Wrap(err, "claiming refund after decryption failure")
			}

			log.Printf("Refund claimed in transaction %x", receipt.TxHash[:])
			return nil

		} else if err != nil {
			log.Fatalf("Error decrypting content: %s", err)
		}
		log.Printf("Complete, decrypted content is in %s", outFileName)
		return nil

	case err := <-subErrChan:
		return errors.Wrap(err, "watching for decryption key event")
	}
}
