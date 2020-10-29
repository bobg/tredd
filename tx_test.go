package tredd

import (
	"context"
	"crypto/sha256"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/common"

	"github.com/bobg/tredd/contract"
	"github.com/bobg/tredd/testutil"
)

var big1 = big.NewInt(1)

func TestProposeCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	_, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	// Canceling before the reveal deadline should fail.
	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err == nil {
		t.Fatal("expected a cancel before the reveal deadline to fail")
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProposePayCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	_, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big1
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// xxx check buyer's balance is decreased

	// Canceling before the reveal deadline should fail.
	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err == nil {
		t.Fatal("expected a cancel before the reveal deadline to fail")
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err != nil {
		t.Fatal(err)
	}

	// xxx check buyer's balance is restored (modulo gas)
}

func TestProposeRevealCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big1
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	con, _, err = RevealKey(ctx, harness.Client, time.Unix(0, 0), &txOpts, contractAddr, testutil.DecryptionKey, common.Address{}, big1, big1, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err == nil {
		t.Fatalf("expected a cancel after the key is revealed to fail")
	}
}

func TestProposeRevealRefundOK(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big1
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// Reveal the wrong key.
	key := testutil.DecryptionKey
	key[0] ^= 1

	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	con, _, err = RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, key, common.Address{}, big1, big1, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs()
	if err != nil {
		t.Fatal(err)
	}

	_, err = ClaimRefund(ctx, harness.Client, harness.Buyer, con, 0, cipherChunk0, clearHash0, cipherProof, clearProof)
	if err != nil {
		t.Fatal(err)
	}

	// xxx check buyer collected payment and collateral
}

func TestProposeRevealRefundFail(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big1
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// Reveal the right key.
	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	con, _, err = RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big1, big1, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs()
	if err != nil {
		t.Fatal(err)
	}

	_, err = ClaimRefund(ctx, harness.Client, harness.Buyer, con, 0, cipherChunk0, clearHash0, cipherProof, clearProof)
	if err == nil {
		t.Fatalf("expected refund attempt to fail after reveal of correct key")
	}
}

func TestProposeRevealClaimPayment(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big1, big1, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big1
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// Reveal the right key.
	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	con, _, err = RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big1, big1, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.AdjustTime(testutil.RefundDeadlineSecs * time.Second)

	_, err = ClaimPayment(ctx, harness.Client, harness.Seller, contractAddr)
	if err != nil {
		t.Fatal(err)
	}

	// xxx check buyer collected payment and collateral
}

func createProofs() (clearHash0 [32]byte, cipherChunk0 []byte, clearProof, cipherProof merkle.Proof, err error) {
	var f io.ReadCloser
	f, err = os.Open("testdata/udhr.txt")
	if err != nil {
		return
	}
	defer f.Close()

	errch := make(chan error, 1)
	pr, pw := io.Pipe()
	go func() {
		defer close(errch)
		defer pw.Close()
		_, err := Serve(pw, f, testutil.DecryptionKey)
		errch <- err
	}()

	var (
		clearMT  *merkle.HTree
		cipherMT *merkle.Tree
	)
	err = Receive(
		pr,
		func(clearHash [32]byte, i uint64) error {
			if i == 0 {
				clearHash0 = clearHash
				clearMT = merkle.NewProofHTree(sha256.New(), clearHash[:])
			}
			clearMT.Add(clearHash[:])
			return nil
		},
		func(cipherChunk []byte, i uint64) error {
			prefixedCipherChunk := PrefixChunk(i, cipherChunk)
			if i == 0 {
				cipherChunk0 = cipherChunk
				cipherMT = merkle.NewProofTree(sha256.New(), prefixedCipherChunk)
			}
			cipherMT.Add(prefixedCipherChunk)
			return nil
		},
	)
	if err != nil {
		return
	}
	err = <-errch
	if err != nil {
		return
	}

	return clearHash0, cipherChunk0, clearMT.Proof(), cipherMT.Proof(), nil
}
