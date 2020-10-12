package tredd

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ProposePayment publishes a new instance of the Tredd contract instantiated with the given parameters.
// It also approves a transfer for `amount` tokens of `tokenType` to the contract
// and then calls the contract's Pay method.
func ProposePayment(
	ctx context.Context,
	client *ethclient.Client, // see ethclient.Dial
	buyer *bind.TransactOpts, // see bind.NewTransactor
	seller common.Address,
	tokenType common.Address,
	amount, collateral *big.Int,
	clearRoot, cipherRoot [32]byte,
	revealDeadline, refundDeadline time.Time,
) (*Tredd, error) {
	token, err := NewERC20(tokenType, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating token")
	}

	contractAddr, deployTx, con, err := DeployTredd(buyer, client, seller, tokenType, amount, collateral, clearRoot, cipherRoot, revealDeadline.Unix(), refundDeadline.Unix())
	if err != nil {
		return nil, errors.Wrap(err, "deploying contract")
	}

	approveTx, err := token.Approve(buyer, contractAddr, amount)
	if err != nil {
		return nil, errors.Wrap(err, "approving token transfer")
	}

	// TODO: double-check that these WaitMined calls are needed before the Pay call.
	g, ctx2 := errgroup.WithContext(ctx)
	g.Go(func() error {
		_, err := bind.WaitMined(ctx2, client, deployTx)
		return err
	})
	g.Go(func() error {
		_, err := bind.WaitMined(ctx2, client, approveTx)
		return err
	})
	err = g.Wait()
	if err != nil {
		return nil, errors.Wrap(err, "waiting for contract deployment and/or transfer approval")
	}

	payTx, err := con.Pay(buyer)
	if err != nil {
		return nil, errors.Wrap(err, "making payment")
	}

	// Wait for payTx to be mined on-chain.
	_, err = bind.WaitMined(ctx, client, payTx)
	return con, errors.Wrap(err, "awaiting payment transaction")
}

// After the reveal deadline, if no reveal has happened, the buyer cancels the contract.
func Cancel(ctx context.Context, client *ethclient.Client, buyer *bind.TransactOpts, contractAddr common.Address) (*types.Receipt, error) {
	con, err := NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}
	return con.CallCancel(ctx, client, buyer)
}

func (con *Tredd) CallCancel(ctx context.Context, client *ethclient.Client, buyer *bind.TransactOpts) (*types.Receipt, error) {
	tx, err := con.Cancel(buyer)
	if err != nil {
		return nil, errors.Wrap(err, "canceling contract")
	}
	return bind.WaitMined(ctx, client, tx)
}

// RevealKey updates a Tredd contract on-chain by adding the decryption key.
// It also approves a collateral transfer.
func RevealKey(
	ctx context.Context,
	client *ethclient.Client, // see ethclient.Dial
	seller *bind.TransactOpts, // see bind.NewTransactor
	contractAddr common.Address,
	key [32]byte,
	wantTokenType common.Address,
	wantAmount, wantCollateral *big.Int,
	wantClearRoot, wantCipherRoot [32]byte,
) (*types.Receipt, error) {
	con, err := NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}

	callOpts := &bind.CallOpts{Context: ctx}

	gotTokenType, err := con.MTokenType(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mTokenType")
	}
	if gotTokenType != wantTokenType {
		return nil, fmt.Errorf("got token type %s, want %s", gotTokenType.Hex(), wantTokenType.Hex())
	}

	gotAmount, err := con.MAmount(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mAmount")
	}
	if gotAmount.Cmp(wantAmount) != 0 {
		return nil, fmt.Errorf("got amount %s, want %s", gotAmount, wantAmount)
	}

	gotCollateral, err := con.MCollateral(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mCollateral")
	}
	if gotCollateral.Cmp(wantCollateral) != 0 {
		return nil, fmt.Errorf("got collateral %s, want %s", gotCollateral, wantCollateral)
	}

	gotCipherRoot, err := con.MCipherRoot(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mCipherRoot")
	}
	if gotCipherRoot != wantCipherRoot {
		return nil, fmt.Errorf("got cipher root %x, want %x", gotCipherRoot[:], wantCipherRoot[:])
	}

	gotClearRoot, err := con.MClearRoot(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mClearRoot")
	}
	if gotClearRoot != wantClearRoot {
		return nil, fmt.Errorf("got clear root %x, want %x", gotClearRoot[:], wantClearRoot[:])
	}

	gotRefundDeadline, err := con.MRefundDeadline(callOpts)
	if err != nil {
		return nil, errors.Wrap(err, "getting mRefundDeadline")
	}
	wantRefundDeadline := time.Now().Add(time.Hour) // TODO: parameterize?
	if gotRefundDeadline > wantRefundDeadline.Unix() {
		return nil, fmt.Errorf("refund deadline is more than %s in the future", time.Hour)
	}

	token, err := NewERC20(wantTokenType, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating token")
	}

	_, err = token.Approve(seller, contractAddr, wantCollateral)
	if err != nil {
		return nil, errors.Wrap(err, "approving token transfer")
	}

	// TODO: Does the approve transaction have to be mined before the reveal transaction will work?

	tx, err := con.Reveal(seller, key)
	if err != nil {
		return nil, errors.Wrap(err, "invoking ClaimPayment")
	}
	return bind.WaitMined(ctx, client, tx)
}

// ClaimPayment constructs a seller-claims-payment transaction,
// rehydrating and invoking a Tredd contract from the utxo state (identified by the information in r).
func ClaimPayment(
	ctx context.Context,
	client *ethclient.Client,
	seller *bind.TransactOpts,
	contractAddr common.Address,
) (*types.Receipt, error) {
	con, err := NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}
	tx, err := con.ClaimPayment(seller)
	if err != nil {
		return nil, errors.Wrap(err, "invoking ClaimPayment")
	}
	return bind.WaitMined(ctx, client, tx)
}

// ClaimRefund constructs a buyer-claims-refund transaction,
// rehydrating a Tredd contract from the utxo state (identified by the information in r)
// and calling it with the necessary proofs and other information.
func ClaimRefund(
	ctx context.Context,
	client *ethclient.Client,
	buyer *bind.TransactOpts,
	contractAddr common.Address,
	index int64,
	cipherChunk []byte,
	clearHash [32]byte,
	cipherProof, clearProof merkle.Proof,
) (*types.Receipt, error) {
	con, err := NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}
	return con.CallRefund(ctx, client, buyer, index, cipherChunk, clearHash, cipherProof, clearProof)
}

func (con *Tredd) CallRefund(
	ctx context.Context,
	client *ethclient.Client,
	buyer *bind.TransactOpts,
	index int64,
	cipherChunk []byte,
	clearHash [32]byte,
	cipherProof, clearProof merkle.Proof,
) (*types.Receipt, error) {
	var (
		treddCipherProof = toTreddProof(cipherProof)
		treddClearProof  = toTreddProof(clearProof)
	)

	tx, err := con.Refund(buyer, uint64(index), cipherChunk, clearHash, treddCipherProof, treddClearProof)
	if err != nil {
		return nil, errors.Wrap(err, "invoking Refund")
	}
	return bind.WaitMined(ctx, client, tx)
}

func toTreddProof(proof merkle.Proof) []TreddProofStep {
	result := make([]TreddProofStep, 0, len(proof))
	for _, step := range proof {
		result = append(result, TreddProofStep{H: step.H, Left: step.Left})
	}
	return result
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

// ParseResult holds the values parsed from a Tredd contract.
// If the transaction is complete
// (i.e., the seller has added the "reveal-key" call),
// all of the fields will be filled in.
// If the transaction is partial, some fields will be uninitialized.
type ParseResult struct {
	ContractAddr common.Address

	// Amount is the amount of the buyer's payment (not including the seller's collateral).
	Amount    int64
	TokenType string

	ClearRoot      [32]byte
	CipherRoot     [32]byte
	RevealDeadline time.Time
	RefundDeadline time.Time
	Buyer          common.Address
	Seller         common.Address
	Key            [32]byte
}
