package tedd

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"i10r.io/crypto/ed25519"
	"i10r.io/protocol/bc"
	"i10r.io/protocol/txbuilder/standard"
	"i10r.io/protocol/txvm/asm"
	"i10r.io/protocol/txvm/op"
	"i10r.io/protocol/txvm/txvmutil"
)

type Payee struct {
	Quorum  int
	Pubkeys []ed25519.PublicKey
	Refdata []byte
	Tags    []byte
}

func BuildPartialPaymentTx(
	ctx context.Context,
	buyer Payee,
	amount int64,
	assetID bc.Hash,
	clearRoot, cipherRoot [32]byte,
	revealDeadline, refundDeadline time.Time,
	reserver Reserver,
) ([]byte, error) {
	reservation, err := reserver.Reserve(ctx, amount, assetID)
	if err != nil {
		// xxx
	}

	// Where the TEDD contract log entries start.
	teddLogPos := int64(len(reservations.UTXOs))
	if reservation.Change > 0 {
		teddLogPos++
	}

	// With the knowledge of the input args and the TEDD log position,
	// construct the signature program for spending these utxos.
	buf := new(bytes.Buffer)
	// Make sure change is sent to the right place.
	if reservation.Change > 0 {
		fmt.Fprintf(buf, "%d peeklog\n", teddLogPos-1)
		// xxx make sure it's {'O', x'0000...', outputID} (compute the right outputID)
	}

	fmt.Fprintf(buf, "%d peeklog untuple\n", teddLogPos)
	fmt.Fprintf(buf, "3 eq verify\n")
	fmt.Fprintf(buf, "2 roll 'L' eq verify\n")
	fmt.Fprintf(buf, "swap x'%x' eq verify\n", xxxTEDDContractSeed)
	fmt.Fprintf(buf, "%d eq verify\n", revealDeadline.Unix())

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+1)
	fmt.Fprintf(buf, "%d eq verify\n", refundDeadline.Unix())
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+2)
	fmt.Fprintf(buf, "x'%x' eq verify\n", cipherRoot[:])
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+3)
	fmt.Fprintf(buf, "x'%x' eq verify\n", clearRoot[:])
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+4)
	// xxx check buyer

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+5)
	fmt.Fprintf(buf, "%d eq verify\n", amount)
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+6)
	fmt.Fprintf(buf, "x'%x' eq verify\n", assetID[:])
	fmt.Fprintf(buf, "drop drop\n")

	sigprog, err := asm.Assemble(buf.String())
	if err != nil {
		// xxx
	}

	anchoredSigprog := make([]byte, 32+len(sigprog))
	copy(anchoredSigprog, sigprog)

	b := new(txvmutil.Builder)
	for i, utxo := range reservation.UTXOs {
		standard.SpendMultisig(b, xxxquorum, xxxpubkeys, utxo.Amount, utxo.AssetID, utxo.Anchor, standard.PayToMultisigSeed2[:])
		// arg stack: [<value> <deferred contract>]
		b.Op(op.Get) // contract stack: [<deferred contract>] arg stack: [<value>]

		copy(anchoredSigprog[len(sigprog):], utxo.Anchor) // this is what to sign
		q := xxxquorum
		for _, pubkey := range xxxpubkeys {
			var sig []byte
			if q > 0 {
				if prv := xxxcansign(pubkey); prv != nil {
					sig = ed25519.Sign(prv, anchoredSigprog)
					q--
				}
			}
			b.PushdataBytes(sig).Op(op.Put)
		}
		if q > 0 {
			// xxx err - too few sigs
		}
		b.PushdataBytes(sigprog).Op(op.Put)
		b.Op(op.Call)
		b.Op(op.Get) // get the value from the arg stack
		if i > 0 {
			b.Op(op.Merge)
		}
	}
	if reservation.Change > 0 {
		b.PushdataInt64(reservation.Change).Op(op.Split)
		// xxx output top value to buyer
	}

	b.PushdataBytes(xxxTEDDContract).Op(op.Contract)

	b.Op(op.Put) // payment, which was already on the contract stack
	// xxx put buyer
	b.PushdataBytes(clearRoot[:]).Op(op.Put)
	b.PushdataBytes(cipherRoot[:]).Op(op.Put)
	b.PushdataInt64(refundDeadline.Unix()).Op(op.Put)
	b.PushdataInt64(revealDeadline.Unix()).Op(op.Put)

	b.Op(op.Call)

	return b.Build()
}
