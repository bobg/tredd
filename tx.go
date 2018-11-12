package tedd

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/bobg/merkle"
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

// ProposePayment constructs a partial transaction in which the buyer commits funds to the tedd contract.
func ProposePayment(
	ctx context.Context,
	buyer ed25519.PublicKey,
	amount int64,
	assetID bc.Hash,
	clearRoot, cipherRoot [32]byte,
	now, revealDeadline, refundDeadline time.Time,
	reserver Reserver,
	signer Signer,
) ([]byte, error) {
	reservation, err := reserver.Reserve(ctx, amount, assetID, now, revealDeadline)
	if err != nil {
		return nil, errors.Wrap(err, "reserving utxos")
	}

	// Where the TEDD contract log entries start.
	utxos := reservation.UTXOs()
	teddLogPos := 2 * int64(len(utxos)) // one 'I' and one 'L' log entry per standard input

	// With the knowledge of the input args and the TEDD log position,
	// construct the signature program for spending these utxos.
	buf := new(bytes.Buffer)

	fmt.Fprint(buf, "[")

	if reservation.Change() > 0 {
		teddLogPos += 3 // one 'O' and two 'L' log entries
		fmt.Fprintf(buf, "%d peeklog untuple\n", teddLogPos-1)

		// Have to make sure this log entry is {'O', seed, outputID}.
		// Computing the right outputID means simulating the merges and the split below to get the change value's anchor.

		var anchor [32]byte
		copy(anchor[:], utxos[0].Anchor())

		for i := 1; i < len(utxos); i++ {
			var inp [64]byte
			copy(inp[:32], utxos[i].Anchor())
			copy(inp[32:], anchor[:])
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
		fmt.Fprintf(buf, "drop\n")
		fmt.Fprintf(buf, "'O' eq verify\n")
	}

	fmt.Fprintf(buf, "%d peeklog untuple\n", teddLogPos)
	fmt.Fprintf(buf, "4 eq verify\n")
	fmt.Fprintf(buf, "3 roll 'R' eq verify\n") // xxx use txvm.TimerangeCode and other such constants
	fmt.Fprintf(buf, "2 roll x'%x' eq verify\n", teddContractSeed[:])
	fmt.Fprintf(buf, "%d eq verify\n", bc.Millis(revealDeadline))
	fmt.Fprintf(buf, "0 eq verify\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+1)
	fmt.Fprintf(buf, "%d eq verify\n", bc.Millis(refundDeadline))
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+2)
	fmt.Fprintf(buf, "x'%x' eq verify\n", buyer)
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+3)
	fmt.Fprintf(buf, "x'%x' eq verify\n", cipherRoot[:])
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+4)
	fmt.Fprintf(buf, "x'%x' eq verify\n", clearRoot[:])
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+5)
	fmt.Fprintf(buf, "%d eq verify\n", amount)
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprintf(buf, "%d peeklog untuple drop\n", teddLogPos+6)
	fmt.Fprintf(buf, "x'%x' eq verify\n", assetID.Bytes())
	fmt.Fprintf(buf, "drop drop\n")

	fmt.Fprint(buf, "] yield")

	sigprog, err := asm.Assemble(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "assembling signature program")
	}

	anchoredSigprog := make([]byte, 32+len(sigprog))
	copy(anchoredSigprog, sigprog)

	b := new(txvmutil.Builder)
	for i, utxo := range reservation.UTXOs() {
		b.PushdataBytes([]byte{}).Op(op.Put)
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
		b.Op(op.Call) // arg stack is again [<value> <deferred contract>]

		// Get the value from the arg stack, leave the deferred contract there.
		b.Op(op.Get).Op(op.Get).PushdataInt64(1).Op(op.Roll).Op(op.Put)

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
		b.PushdataBytes(standard.PayToMultisigProg2).Op(op.Contract).Op(op.Call)
	}

	b.PushdataBytes(teddContractProg).Op(op.Contract)
	b.PushdataInt64(1).Op(op.Roll)

	b.Op(op.Put) // payment, which was already on the contract stack
	b.PushdataBytes(clearRoot[:]).Op(op.Put)
	b.PushdataBytes(cipherRoot[:]).Op(op.Put)
	b.PushdataBytes(buyer).Op(op.Put)
	b.PushdataInt64(int64(bc.Millis(refundDeadline))).Op(op.Put) // xxx range checking
	b.PushdataInt64(int64(bc.Millis(revealDeadline))).Op(op.Put) // xxx range checking

	b.Op(op.Call)

	// con stack is now empty
	// arg stack is sigprog sigprog ... teddcontract (all deferred)

	b.Op(op.Get) // move tedd contract back to con stack

	// Now that the first phase of the tedd contract has run and begun to populate the tx log,
	// the sig progs, which check the log, can run.
	for i := 0; i < len(reservation.UTXOs()); i++ {
		b.Op(op.Get).Op(op.Call)
	}

	return b.Build(), nil
}

// RevealKey completes the partial transaction in paymentProposal (which came from ProposePayment).
// The tedd contract is on the con stack. The arg stack is empty.
func RevealKey(
	ctx context.Context,
	paymentProposal []byte,
	seller ed25519.PublicKey,
	key [32]byte,
	amount int64,
	assetID bc.Hash,
	reserver Reserver,
	signer Signer,
	wantClearRoot, wantCipherRoot [32]byte,
	now, wantRevealDeadline, wantRefundDeadline time.Time,
) ([]byte, error) {
	parsed := ParseLog(paymentProposal)
	if parsed == nil {
		return nil, errors.New("could not parse payment proposal")
	}
	if parsed.RevealDeadline.Unix() != wantRevealDeadline.Unix() {
		return nil, fmt.Errorf("got reveal deadline %s, want %s", parsed.RevealDeadline, wantRevealDeadline)
	}
	if parsed.RefundDeadline.Unix() != wantRefundDeadline.Unix() {
		return nil, fmt.Errorf("got refund deadline %s, want %s", parsed.RefundDeadline, wantRefundDeadline)
	}
	if !bytes.Equal(parsed.CipherRoot, wantCipherRoot[:]) {
		return nil, fmt.Errorf("got cipher root %x, want %x", parsed.CipherRoot, wantCipherRoot[:])
	}
	if !bytes.Equal(parsed.ClearRoot, wantClearRoot[:]) {
		return nil, fmt.Errorf("got clear root %x, want %x", parsed.ClearRoot, wantClearRoot[:])
	}
	if parsed.Amount != amount {
		return nil, fmt.Errorf("got amount %d, want %d", parsed.Amount, amount)
	}
	if !bytes.Equal(parsed.AssetID, assetID.Bytes()) {
		return nil, fmt.Errorf("got asset ID %x, want %x", parsed.AssetID, assetID.Bytes())
	}

	reservation, err := reserver.Reserve(ctx, amount, assetID, now, wantRevealDeadline)
	if err != nil {
		return nil, errors.Wrap(err, "reserving utxos")
	}

	b := new(txvmutil.Builder)

	for i, utxo := range reservation.UTXOs() {
		b.PushdataBytes([]byte{}).Op(op.Put)
		standard.SpendMultisig(b, 1, []ed25519.PublicKey{seller}, utxo.Amount(), utxo.AssetID(), utxo.Anchor(), standard.PayToMultisigSeed2[:])
		// arg stack: [<value> <deferred contract>]
		b.Op(op.Get).Op(op.Get)
		b.PushdataInt64(1).Op(op.Roll)
		b.Op(op.Put)
		if i > 0 {
			b.Op(op.Merge)
		}
	}
	// con stack: teddcontract collateral
	// arg stack: sigcheck sigcheck ...
	if reservation.Change() > 0 {
		b.PushdataInt64(reservation.Change()).Op(op.Split)
		// con stack: teddcontract collateral change
		b.PushdataBytes([]byte{}).Op(op.Put)
		b.PushdataBytes([]byte{}).Op(op.Put)
		b.Op(op.Put)
		b.PushdataBytes(seller).PushdataInt64(1).Op(op.Tuple).Op(op.Put)
		b.PushdataInt64(1).Op(op.Put)
		b.PushdataBytes(standard.PayToMultisigProg2).Op(op.Contract).Op(op.Call)
	}

	spendProg := b.Build()

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "x'%x' exec\n", spendProg)

	fmt.Fprintf(buf, "%d split\n", amount) // con stack: teddcontract zeroval collateral

	fmt.Fprintf(buf, "x'%x' put\n", seller)
	fmt.Fprintf(buf, "x'%x' put\n", key[:])

	fmt.Fprintf(buf, "put\n")  // move collateral to arg stack
	fmt.Fprintf(buf, "swap\n") // con stack: zeroval teddcontract
	fmt.Fprintf(buf, "call\n") // con stack: zeroval
	fmt.Fprintf(buf, "finalize\n")

	tx1, err := asm.Assemble(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "assembling unsigned transaction")
	}
	tx1 = append(paymentProposal, tx1...)

	vm, err := txvm.Validate(tx1, 3, math.MaxInt64, txvm.StopAfterFinalize)
	if err != nil {
		return nil, errors.Wrap(err, "computing transaction ID")
	}

	// sign seller utxos
	buf = new(bytes.Buffer)
	sigprog := standard.VerifyTxID(vm.TxID)
	for i := len(reservation.UTXOs()) - 1; i >= 0; i-- {
		utxo := reservation.UTXOs()[i]
		anchoredSigprog := append([]byte{}, sigprog...)
		anchoredSigprog = append(anchoredSigprog, utxo.Anchor()...)
		sig, err := signer(anchoredSigprog)
		if err != nil {
			return nil, errors.Wrap(err, "computing signature")
		}
		fmt.Fprintf(buf, "get x'%x' put x'%x' put call\n", sig, sigprog)
	}
	tx2, err := asm.Assemble(buf.String())
	if err != nil {
		return nil, errors.Wrap(err, "assembling signature section")
	}
	return append(tx1, tx2...), nil
}

type Redeem struct {
	RefundDeadline        time.Time
	Buyer, Seller         ed25519.PublicKey
	Amount                int64 // sum of buyer payment + collateral
	AssetID               bc.Hash
	Anchor                [32]byte
	CipherRoot, ClearRoot [32]byte
	Key                   [32]byte
}

func redeem(r *Redeem) *bytes.Buffer {
	buf := new(bytes.Buffer)
	fmt.Fprintf(
		buf,
		"{'C', x'%x', x'%x', {'Z', %d}, {'S', x'%x'}, {'V', %d, x'%x', x'%x'}, {'S', x'%x'}, {'S', x'%x'}, {'S', x'%x'}, {'S', x'%x'}} input\n",
		teddContractSeed,
		redemptionProg,
		bc.Millis(r.RefundDeadline),
		r.Buyer,
		r.Amount,
		r.AssetID.Bytes(),
		r.Anchor[:],
		r.CipherRoot[:],
		r.ClearRoot[:],
		r.Key[:],
		r.Seller,
	)
	return buf
}

func ClaimPayment(r *Redeem) ([]byte, error) {
	buf := redeem(r)
	fmt.Fprintln(buf, "0 put call")
	fmt.Fprintln(buf, "get finalize")
	return asm.Assemble(buf.String())
}

func ClaimRefund(r *Redeem, index int64, cipherChunk []byte, clearHash []byte, cipherProof, clearProof merkle.Proof) ([]byte, error) {
	var prefix [binary.MaxVarintLen64]byte
	m := binary.PutUvarint(prefix[:], uint64(index))

	buf := redeem(r)
	renderProof(buf, cipherProof)
	fmt.Fprintln(buf, "put")
	renderProof(buf, clearProof)
	fmt.Fprintln(buf, "put")
	fmt.Fprintf(buf, "x'%x' put\n", clearHash)
	fmt.Fprintf(buf, "x'%x' put\n", cipherChunk)
	fmt.Fprintf(buf, "x'%x' put\n", prefix[:m])
	fmt.Fprintln(buf, "1 put call")
	fmt.Fprintln(buf, "get finalize")
	return asm.Assemble(buf.String())
}

func renderProof(w io.Writer, proof merkle.Proof) {
	fmt.Fprint(w, "{")
	for i := len(proof) - 1; i >= 0; i-- {
		if i < len(proof)-1 {
			fmt.Fprint(w, ", ")
		}
		var isLeft int64
		if proof[i].Left {
			isLeft = 1
		}
		fmt.Fprintf(w, "x'%x', %d", proof[i].H, isLeft)
	}
	fmt.Fprintln(w, "}")
}

type ParseResult struct {
	Amount         int64
	AssetID        []byte
	Anchor1        []byte
	Anchor2        []byte
	ClearRoot      []byte
	CipherRoot     []byte
	RevealDeadline time.Time
	RefundDeadline time.Time
	Buyer          ed25519.PublicKey
	Seller         ed25519.PublicKey
	Key            []byte
	OutputID       []byte
}

func ParseLog(prog []byte) *ParseResult {
	vm, err := txvm.Validate(prog, 3, math.MaxInt64, txvm.StopAfterFinalize)
	if vm == nil || err != nil {
		return nil
	}
	var res *ParseResult
	for i, item := range vm.Log {
		if len(item) != 4 {
			continue
		}
		code, ok := item[0].(txvm.Bytes)
		if !ok {
			continue
		}
		if !bytes.Equal(code, []byte{'R'}) {
			continue
		}
		if !bytes.Equal(item[1].(txvm.Bytes), teddContractSeed[:]) {
			continue
		}
		res = &ParseResult{
			RevealDeadline: bc.FromMillis(uint64(item[3].(txvm.Int))), // xxx range checking
			RefundDeadline: bc.FromMillis(uint64(vm.Log[i+1][2].(txvm.Int))),
			Buyer:          ed25519.PublicKey(vm.Log[i+2][2].(txvm.Bytes)),
			CipherRoot:     vm.Log[i+3][2].(txvm.Bytes),
			ClearRoot:      vm.Log[i+4][2].(txvm.Bytes),
			Amount:         int64(vm.Log[i+5][2].(txvm.Int)),
			AssetID:        vm.Log[i+6][2].(txvm.Bytes),
			Anchor1:        vm.Log[i+7][2].(txvm.Bytes),
		}
		for j := i + 1; j < len(vm.Log); j++ {
			item := vm.Log[j]
			if len(item) != 3 {
				continue
			}
			if !bytes.Equal(code, []byte{'L'}) {
				continue
			}
			if !bytes.Equal(item[1].(txvm.Bytes), teddContractSeed[:]) {
				continue
			}
			res.Anchor1 = vm.Log[j][2].(txvm.Bytes)
			res.Key = vm.Log[j+1][2].(txvm.Bytes)
			res.Seller = ed25519.PublicKey(vm.Log[j+2][2].(txvm.Bytes))
			res.OutputID = vm.Log[j+3][2].(txvm.Bytes)
			break
		}
		break
	}
	return res
}
