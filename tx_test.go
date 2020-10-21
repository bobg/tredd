package tredd

import (
	"context"
	"os"
	"testing"

	"github.com/bobg/tredd/testutil"
)

func TestTx(t *testing.T) {
	harness, err := testutil.NewHarness()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	f, err := os.Open("testdata/udhr.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// test ProposePayment/Cancel
	err = harness.Deploy(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// test ProposePayment/Pay/Cancel
	// test ProposePayment/RevealKey/Cancel (should fail)
	// test ProposePayment/RevealKey/ClaimRefund with happy values (should fail)
	// test ProposePayment/RevealKey/ClaimRefund with sad values (should succeed)
	// test ProposePayment/RevealKey/ClaimPayment
}
