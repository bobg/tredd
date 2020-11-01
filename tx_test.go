package tredd

import (
	"context"
	"crypto/sha256"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/bobg/merkle/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bobg/tredd/testutil"
)

var (
	big2 = big.NewInt(2)
	big3 = big.NewInt(3)
)

func TestProposeCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	// Canceling before the reveal deadline should fail.
	_, err = Cancel(ctx, harness.Client, harness.Buyer, con)
	if err == nil {
		t.Fatal("expected a cancel before the reveal deadline to fail")
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	rcpt, err := Cancel(ctx, harness.Client, harness.Buyer, con)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance += 3
	harness.BuyerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
		contractBal, err := harness.Client.BalanceAt(ctx, contractAddr, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("contract balance is %s", contractBal)
	}
}

func TestProposeRevealCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	con, rcpt, err := RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big3, big2, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance -= 2
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
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

	contractAddr, con, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	// Reveal the wrong key.
	key := testutil.DecryptionKey
	key[0] ^= 1

	con, rcpt, err := RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, key, common.Address{}, big3, big2, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance -= 2
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs(false)
	if err != nil {
		t.Fatal(err)
	}

	rcpt, err = ClaimRefund(ctx, harness.Client, harness.Buyer, con, 0, cipherChunk0, clearHash0, cipherProof, clearProof)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance += 5
	harness.BuyerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestProposeRevealRefundFail(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	// Reveal the right key.
	con, rcpt, err := RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big3, big2, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance -= 2
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs(false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ClaimRefund(ctx, harness.Client, harness.Buyer, con, 0, cipherChunk0, clearHash0, cipherProof, clearProof)
	if err == nil {
		t.Fatalf("expected refund attempt to fail after reveal of correct key")
	}
}

func TestProposeRevealRefundFraud(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, con, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	// Reveal the right key.
	con, rcpt, err := RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big3, big2, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance -= 2
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs(true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ClaimRefund(ctx, harness.Client, harness.Buyer, con, 0, cipherChunk0, clearHash0, cipherProof, clearProof)
	if err == nil {
		t.Fatalf("expected refund attempt to fail after reveal of correct key")
	} else {
		t.Log(err)
	}
}

func TestProposeRevealClaimPayment(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	contractAddr, _, rcpts, err := ProposePayment(ctx, harness.Client, harness.Buyer, harness.Seller.From, common.Address{}, big3, big2, testutil.ClearRoot, testutil.CipherRoot, harness.RevealDeadline, harness.RefundDeadline)
	if err != nil {
		t.Fatal(err)
	}

	harness.BuyerBalance -= 3
	harness.BuyerBalance -= gasUsed(rcpts)
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	// Reveal the right key.
	_, rcpt, err := RevealKey(ctx, harness.Client, time.Unix(0, 0), harness.Seller, contractAddr, testutil.DecryptionKey, common.Address{}, big3, big2, harness.RevealDeadline, harness.RefundDeadline, testutil.ClearRoot, testutil.CipherRoot)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance -= 2
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
	}

	harness.Client.AdjustTime(testutil.RefundDeadlineSecs * time.Second)

	rcpt, err = ClaimPayment(ctx, harness.Client, harness.Seller, contractAddr)
	if err != nil {
		t.Fatal(err)
	}

	harness.SellerBalance += 5
	harness.SellerBalance -= rcpt.GasUsed
	err = harness.CheckBalances(ctx)
	if err != nil {
		t.Error(err)
		contractBal, err := harness.Client.BalanceAt(ctx, contractAddr, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("contract balance is %s", contractBal)
	}
}

func createProofs(fraud bool) (clearHash0 [32]byte, cipherChunk0 []byte, clearProof, cipherProof merkle.Proof, err error) {
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
		clearMT  *merkle.Tree
		cipherMT *merkle.Tree
	)
	err = Receive(
		pr,
		func(clearHash [32]byte, i uint64) error {
			prefixedClearHash := Prefix(i, clearHash[:])
			if i == 0 {
				clearHash0 = clearHash
				clearMT = merkle.NewProofTree(sha256.New(), prefixedClearHash)
			} else if i == 1 && fraud {
				clearMT = merkle.NewProofTree(sha256.New(), prefixedClearHash)
				clearHash0Copy := clearHash0 // the merkle tree takes ownership of this memory
				clearMT.Add(Prefix(0, clearHash0Copy[:]))
				clearHash0 = clearHash
			}
			clearMT.Add(prefixedClearHash)
			return nil
		},
		func(cipherChunk []byte, i uint64) error {
			prefixedCipherChunk := Prefix(i, cipherChunk)
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

func gasUsed(rcpts []*types.Receipt) uint64 {
	var g uint64
	for _, r := range rcpts {
		g += r.GasUsed
	}
	return g
}
