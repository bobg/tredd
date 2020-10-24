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
	"github.com/pkg/errors"

	"github.com/bobg/tredd/contract"
	"github.com/bobg/tredd/testutil"
)

func TestProposeCancel(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Canceling before the reveal deadline should fail.
	err = harness.Cancel(ctx)
	if err == nil {
		t.Fatal("expected a cancel before the reveal deadline to fail")
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	err = harness.Cancel(ctx)
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

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big.NewInt(1)
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// xxx check buyer's balance is decreased

	// Canceling before the reveal deadline should fail.
	err = harness.Cancel(ctx)
	if err == nil {
		t.Fatal("expected a cancel before the reveal deadline to fail")
	}

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	err = harness.Cancel(ctx)
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

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big.NewInt(1)
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	_, err = con.Reveal(&txOpts, testutil.DecryptionKey)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	harness.Client.AdjustTime(testutil.RevealDeadlineSecs * time.Second)

	err = harness.Cancel(ctx)
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

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big.NewInt(1)
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
	_, err = con.Reveal(&txOpts, key)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs()
	if err != nil {
		t.Fatal(err)
	}

	_, err = con.Refund(harness.Buyer, 0, cipherChunk0, clearHash0, contract.Proof(cipherProof), contract.Proof(clearProof))
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// xxx check buyer collected payment and collateral
}

func TestProposeRevealRefundFail(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big.NewInt(1)
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// Reveal the right key.
	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	_, err = con.Reveal(&txOpts, testutil.DecryptionKey)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	clearHash0, cipherChunk0, clearProof, cipherProof, err := createProofs()
	if err != nil {
		t.Fatal(err)
	}

	_, err = con.Refund(harness.Buyer, 0, cipherChunk0, clearHash0, contract.Proof(cipherProof), contract.Proof(clearProof))
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

	// test ProposePayment/Cancel
	err = harness.ProposePayment(ctx)
	if err != nil {
		t.Fatal(err)
	}

	con, err := harness.Contract()
	if err != nil {
		t.Fatal(err)
	}

	txOpts := *harness.Buyer
	txOpts.Value = big.NewInt(1)
	raw := &contract.TreddRaw{Contract: con}

	_, err = raw.Transfer(&txOpts)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	// Reveal the right key.
	txOpts = *harness.Seller
	txOpts.Value = big.NewInt(1)
	_, err = con.Reveal(&txOpts, testutil.DecryptionKey)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

	harness.Client.AdjustTime(testutil.RefundDeadlineSecs)

	_, err = con.ClaimPayment(harness.Seller)
	if err != nil {
		t.Fatal(err)
	}

	harness.Client.Commit()

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
		_, err = Serve(pw, f, testutil.DecryptionKey)
		errch <- err
	}()

	var (
		wasPartial        = false
		clearMT, cipherMT *merkle.Tree
	)
	for i := uint64(0); ; i++ {
		var clearHash [32]byte
		_, err = io.ReadFull(pr, clearHash[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		var (
			cipherChunk [ChunkSize]byte
			n           int
		)
		n, err = io.ReadFull(pr, cipherChunk[:])
		if err == io.ErrUnexpectedEOF {
			if wasPartial {
				err = errors.New("two partial chunks seen")
				return
			}
			wasPartial = true
		} else if err != nil {
			return
		}

		var (
			prefixedClearHash   = PrefixHash(i, clearHash)
			prefixedCipherChunk = PrefixChunk(i, cipherChunk[:n])
		)

		if i == 0 {
			clearHash0 = clearHash
			cipherChunk0 = cipherChunk[:n]

			clearMT = merkle.NewProofTree(sha256.New(), prefixedClearHash)
			cipherMT = merkle.NewProofTree(sha256.New(), prefixedCipherChunk)
		}

		clearMT.Add(prefixedClearHash)
		cipherMT.Add(prefixedCipherChunk)
	}

	err = <-errch
	if err != nil {
		return
	}

	return clearHash0, cipherChunk0, clearMT.Proof(), cipherMT.Proof(), nil
}
