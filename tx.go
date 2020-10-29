package tredd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/bobg/tredd/contract"
)

type clientType interface {
	bind.ContractBackend
	bind.DeployBackend
}

// ProposePayment publishes a new instance of the Tredd contract instantiated with the given parameters.
// It also approves a transfer for `amount` tokens of `tokenType` to the contract
// and then calls the contract's Pay method.
func ProposePayment(
	ctx context.Context,
	client clientType,
	buyer *bind.TransactOpts,
	seller common.Address,
	tokenType common.Address,
	amount, collateral *big.Int,
	clearRoot, cipherRoot [32]byte,
	revealDeadline, refundDeadline time.Time,
) (common.Address, *contract.Tredd, error) {
	contractAddr, deployTx, con, err := contract.DeployTredd(buyer, client, seller, tokenType, amount, collateral, clearRoot, cipherRoot, revealDeadline.Unix(), refundDeadline.Unix())
	if err != nil {
		return common.Address{}, nil, errors.Wrap(err, "deploying contract")
	}

	_, err = waitMined(ctx, client, deployTx)
	if err != nil {
		return common.Address{}, nil, errors.Wrap(err, "waiting for contract deployment")
	}

	var payTx *types.Transaction
	if IsETH(tokenType) {
		buyer := *buyer
		buyer.Value = amount

		raw := &contract.TreddRaw{Contract: con}
		payTx, err = raw.Transfer(&buyer)
		if err != nil {
			return common.Address{}, nil, errors.Wrap(err, "making payment")
		}
	} else {
		token, err := contract.NewERC20(tokenType, client)
		if err != nil {
			return common.Address{}, nil, errors.Wrap(err, "instantiating token")
		}
		payTx, err = token.Transfer(buyer, contractAddr, amount)
		if err != nil {
			return common.Address{}, nil, errors.Wrap(err, "making payment")
		}
	}

	// Wait for payTx to be mined on-chain.
	_, err = waitMined(ctx, client, payTx)
	return contractAddr, con, errors.Wrap(err, "awaiting payment transaction")
}

// After the reveal deadline, if no reveal has happened, the buyer cancels the contract.
func Cancel(ctx context.Context, client clientType, buyer *bind.TransactOpts, con *contract.Tredd) (*types.Receipt, error) {
	tx, err := con.Cancel(buyer)
	if err != nil {
		return nil, errors.Wrap(err, "canceling contract")
	}
	return waitMined(ctx, client, tx)
}

// The reveal deadline must still be this far in the future when RevealKey is called.
const minRevealDur = 5 * time.Minute

// RevealKey updates a Tredd contract on-chain by adding the decryption key.
// It also approves a collateral transfer.
func RevealKey(
	ctx context.Context,
	client clientType,
	now time.Time,
	seller *bind.TransactOpts,
	contractAddr common.Address,
	key [32]byte,
	wantTokenType common.Address,
	wantAmount, wantCollateral *big.Int,
	wantRevealDeadline, wantRefundDeadline time.Time,
	wantClearRoot, wantCipherRoot [32]byte,
) (*contract.Tredd, *types.Receipt, error) {
	con, err := contract.NewTredd(contractAddr, client)
	if err != nil {
		return nil, nil, errors.Wrap(err, "instantiating deployed contract")
	}

	callOpts := &bind.CallOpts{Context: ctx}

	gotTokenType, err := con.MTokenType(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mTokenType")
	}
	if gotTokenType != wantTokenType {
		return nil, nil, fmt.Errorf("got token type %s, want %s", gotTokenType.Hex(), wantTokenType.Hex())
	}

	gotAmount, err := con.MAmount(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mAmount")
	}
	if gotAmount.Cmp(wantAmount) != 0 {
		return nil, nil, fmt.Errorf("got amount %s, want %s", gotAmount, wantAmount)
	}

	gotCollateral, err := con.MCollateral(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mCollateral")
	}
	if gotCollateral.Cmp(wantCollateral) != 0 {
		return nil, nil, fmt.Errorf("got collateral %s, want %s", gotCollateral, wantCollateral)
	}

	gotClearRoot, err := con.MClearRoot(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mClearRoot")
	}
	if gotClearRoot != wantClearRoot {
		return nil, nil, fmt.Errorf("got clear root %x, want %x", gotClearRoot[:], wantClearRoot[:])
	}

	gotCipherRoot, err := con.MCipherRoot(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mCipherRoot")
	}
	if gotCipherRoot != wantCipherRoot {
		return nil, nil, fmt.Errorf("got cipher root %x, want %x", gotCipherRoot[:], wantCipherRoot[:])
	}

	gotRevealDeadlineSecs, err := con.MRevealDeadline(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mRevealDeadline")
	}
	gotRevealDeadline := time.Unix(gotRevealDeadlineSecs, 0)
	if gotRevealDeadlineSecs != wantRevealDeadline.Unix() { // lop off fractional seconds from wantRevealDeadline
		return nil, nil, fmt.Errorf("reveal deadline is %s, want %s", gotRevealDeadline, wantRevealDeadline)
	}

	if gotRevealDeadline.Sub(now) < minRevealDur {
		return nil, nil, fmt.Errorf("reveal deadline of %s is too soon, or in the past", gotRevealDeadline)
	}

	gotRefundDeadlineSecs, err := con.MRefundDeadline(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mRefundDeadline")
	}
	if gotRefundDeadlineSecs != wantRefundDeadline.Unix() { // lop off fractional seconds from wantRefundDeadline
		return nil, nil, fmt.Errorf("refund deadline is %s, want %s", time.Unix(gotRefundDeadlineSecs, 0), wantRefundDeadline)
	}

	paidAmount, err := con.Paid(callOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "checking paid amount")
	}
	if paidAmount.Cmp(wantAmount) < 0 {
		return nil, nil, fmt.Errorf("contract balance is %s, want %s", paidAmount, wantAmount)
	}

	if !IsETH(wantTokenType) {
		token, err := contract.NewERC20(wantTokenType, client)
		if err != nil {
			return nil, nil, errors.Wrap(err, "instantiating token")
		}
		_, err = token.Approve(seller, contractAddr, wantCollateral)
		if err != nil {
			return nil, nil, errors.Wrap(err, "approving token transfer")
		}
		// TODO: Does the approve transaction have to be mined before the reveal transaction will work?
	}

	revealTxOpts := seller
	if IsETH(wantTokenType) {
		seller := *seller
		seller.Value = wantCollateral
		revealTxOpts = &seller
	}

	revealTx, err := con.Reveal(revealTxOpts, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invoking ClaimPayment")
	}

	receipt, err := waitMined(ctx, client, revealTx)
	return con, receipt, errors.Wrap(err, "waiting for reveal tx to be mined")
}

// ClaimPayment constructs a seller-claims-payment transaction,
// rehydrating and invoking a Tredd contract from the utxo state (identified by the information in r).
func ClaimPayment(
	ctx context.Context,
	client clientType,
	seller *bind.TransactOpts,
	contractAddr common.Address,
) (*types.Receipt, error) {
	con, err := contract.NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}
	tx, err := con.ClaimPayment(seller)
	if err != nil {
		return nil, errors.Wrap(err, "invoking ClaimPayment")
	}
	return waitMined(ctx, client, tx)
}

// ClaimRefund constructs a buyer-claims-refund transaction,
// rehydrating a Tredd contract from the utxo state (identified by the information in r)
// and calling it with the necessary proofs and other information.
func ClaimRefund(
	ctx context.Context,
	client clientType,
	buyer *bind.TransactOpts,
	con *contract.Tredd,
	index uint64,
	cipherChunk []byte,
	clearHash [32]byte,
	cipherProof, clearProof merkle.Proof,
) (*types.Receipt, error) {
	var (
		treddCipherProof = contract.Proof(cipherProof)
		treddClearProof  = contract.Proof(clearProof)
	)

	tx, err := con.Refund(buyer, index, cipherChunk, clearHash, treddCipherProof, treddClearProof)
	if err != nil {
		return nil, errors.Wrap(err, "invoking Refund")
	}
	return waitMined(ctx, client, tx)
}

func waitMined(ctx context.Context, client clientType, tx *types.Transaction) (*types.Receipt, error) {
	if simulated, ok := client.(*backends.SimulatedBackend); ok {
		simulated.Commit()
	}
	return bind.WaitMined(ctx, client, tx)
}

var ethAddr common.Address

func IsETH(tokenType common.Address) bool {
	return tokenType == ethAddr
}
