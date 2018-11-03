package tedd

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txbuilder/standard"
	"github.com/chain/txvm/protocol/txvm"
	"github.com/chain/txvm/protocol/txvm/asm"
	"github.com/chain/txvm/protocol/txvm/op"
	"github.com/chain/txvm/protocol/txvm/txvmutil"
)

type Signer func([]byte) ([]byte, error)

func BuildPartialPaymentTx(
	ctx context.Context,
	buyer ed25519.PublicKey,
	amount int64,
	assetID bc.Hash,
	clearRoot, cipherRoot [32]byte,
	revealDeadline, refundDeadline time.Time,
	reserver Reserver,
	signer Signer,
) ([]byte, error) {
	reservation, err := reserver.Reserve(ctx, amount, assetID, revealDeadline)
	if err != nil {
		return nil, errors.Wrap(err, "reserving utxos")
	}

	// Where the TEDD contract log entries start.
	utxos := reservation.UTXOs()
	teddLogPos := 2 * int64(len(utxos)) // one 'I' and one 'L' log entry per standard input

	// With the knowledge of the input args and the TEDD log position,
	// construct the signature program for spending these utxos.
	buf := new(bytes.Buffer)

	if reservation.Change() > 0 {
		teddLogPos += 3 // one 'O' and two 'L' log entries
		fmt.Fprintf(buf, "%d peeklog untuple\n", teddLogPos-1)

		// Have to make sure this log entry is {'O', seed, outputID}.
		// Computing the right outputID means simulating the merges and the split below to get the change value's anchor.

		var anchor [32]byte
		copy(anchor[:], utxos[0].Anchor())

		for i := 1; i < len(utxos); i++ {
			var inp [64]byte
			copy(inp[:32], anchor[:])
			copy(inp[32:], utxos[i].Anchor())
			anchor = txvm.VMHash("Merge", inp[:])
		}

		anchor = txvm.VMHash("Split2", anchor[:])

		b := new(txvmutil.Builder)
		standard.SpendMultisig(b, 1, []ed25519.PublicKey{buyer}, reservation.Change(), assetID, anchor[:], standard.PayToMultisigSeed2[:])
		snapshot := b.Build()

		// This lops off the "input" and "call" opcodes at the end of standard.SpendMultisig.
		// TODO: refactor SpendMultisig to make the snapshot tuple available separately.
		snapshot = snapshot[:len(snapshot)-2]
		outputID := txvm.VMHash("SnapshotID", snapshot)

		fmt.Fprintf(buf, "3 eq verify\n")
		fmt.Fprintf(buf, "x'%x' eq verify\n", outputID[:])
		fmt.Fprintf(buf, "x'%x' eq verify\n", standard.PayToMultisigSeed2[:])
		fmt.Fprintf(buf, "'O' eq verify\n")
	}

	fmt.Fprintf(buf, "%d peeklog untuple\n", teddLogPos)
	fmt.Fprintf(buf, "3 eq verify\n")
	fmt.Fprintf(buf, "2 roll 'L' eq verify\n")
	fmt.Fprintf(buf, "swap x'%x' eq verify\n", []byte{}) // xxx TEDDContractSeed
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
	fmt.Fprintf(buf, "x'%x' eq verify\n", buyer)
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+5)
	fmt.Fprintf(buf, "%d eq verify\n", amount)
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+6)
	fmt.Fprintf(buf, "x'%x' eq verify\n", assetID.Bytes())
	fmt.Fprintf(buf, "drop drop\n")

	sigprog, err := asm.Assemble(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "assembling signature program")
	}

	anchoredSigprog := make([]byte, 32+len(sigprog))
	copy(anchoredSigprog, sigprog)

	b := new(txvmutil.Builder)
	for i, utxo := range reservation.UTXOs() {
		standard.SpendMultisig(b, 1, []ed25519.PublicKey{buyer}, utxo.Amount(), utxo.AssetID(), utxo.Anchor(), standard.PayToMultisigSeed2[:])
		// arg stack: [<value> <deferred contract>]
		b.Op(op.Get) // contract stack: [<deferred contract>] arg stack: [<value>]

		copy(anchoredSigprog[len(sigprog):], utxo.Anchor()) // this is what to sign
		sig, err := signer(anchoredSigprog)
		if err != nil {
			return nil, errors.Wrap(err, "signing input")
		}
		b.PushdataBytes(sig).Op(op.Put)
		b.PushdataBytes(sigprog).Op(op.Put)
		b.Op(op.Call)

		b.Op(op.Get) // get the value from the arg stack
		if i > 0 {
			b.Op(op.Merge)
		}
	}
	if reservation.Change() > 0 {
		b.PushdataInt64(reservation.Change()).Op(op.Split)

		b.PushdataBytes(nil).Op(op.Put)
		b.PushdataBytes(nil).Op(op.Put)
		b.Op(op.Put)
		b.PushdataBytes(buyer).PushdataInt64(1).Op(op.Tuple).Op(op.Put)
		b.PushdataInt64(1).Op(op.Put)
		b.Concat(standard.PayToMultisigProg2).Op(op.Contract).Op(op.Call)
	}

	b.PushdataBytes(nil /* xxx TEDDContract */).Op(op.Contract)

	b.Op(op.Put) // payment, which was already on the contract stack
	b.PushdataBytes(buyer).Op(op.Put)
	b.PushdataBytes(clearRoot[:]).Op(op.Put)
	b.PushdataBytes(cipherRoot[:]).Op(op.Put)
	b.PushdataInt64(refundDeadline.Unix()).Op(op.Put)
	b.PushdataInt64(revealDeadline.Unix()).Op(op.Put)

	b.Op(op.Call)

	return b.Build(), nil
}
