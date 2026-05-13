package tredd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bobg/errors"
	"github.com/bobg/merkle/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bobg/tredd/contract"
)

var (
	treddABI = contract.NewTredd()
	erc20ABI = contract.NewERC20()
)

type clientType interface {
	bind.ContractBackend
	bind.DeployBackend
}

// committer is implemented by simulated backends that can mine a block on demand.
type committer interface {
	Commit() common.Hash
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
) (common.Address, *bind.BoundContract, []*types.Receipt, error) {
	var rcpts []*types.Receipt

	txOpts := buyer
	if IsETH(tokenType) {
		o := *txOpts
		o.Value = amount
		txOpts = &o
	}

	constructorInput := treddABI.PackConstructor(seller, tokenType, amount, collateral, clearRoot, cipherRoot, uint64(revealDeadline.Unix()), uint64(refundDeadline.Unix()))
	contractAddr, deployTx, err := bind.DeployContract(txOpts, common.FromHex(contract.TreddMetaData.Bin), client, constructorInput)
	if err != nil {
		return common.Address{}, nil, nil, errors.Wrap(err, "deploying contract")
	}

	rcpt, err := waitMined(ctx, client, deployTx)
	if err != nil {
		return common.Address{}, nil, nil, errors.Wrap(err, "waiting for contract deployment")
	}
	rcpts = append(rcpts, rcpt)

	con := treddABI.Instance(client, contractAddr)

	if !IsETH(tokenType) {
		tokenInstance := erc20ABI.Instance(client, tokenType)
		payTx, err := bind.Transact(tokenInstance, buyer, erc20ABI.PackTransfer(contractAddr, amount))
		if err != nil {
			return common.Address{}, nil, nil, errors.Wrap(err, "making payment")
		}
		rcpt, err = waitMined(ctx, client, payTx)
		if err != nil {
			return common.Address{}, nil, nil, errors.Wrap(err, "funding contract")
		}
		rcpts = append(rcpts, rcpt)
	}

	// Wait for payTx to be mined on-chain.
	return contractAddr, con, rcpts, nil
}

// Cancel cancels the contract if, after the reveal deadline, no reveal has happened.
func Cancel(ctx context.Context, client clientType, buyer *bind.TransactOpts, con *bind.BoundContract) (*types.Receipt, error) {
	tx, err := bind.Transact(con, buyer, treddABI.PackCancel())
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
) (*bind.BoundContract, *types.Receipt, error) {
	con := contract.NewInstance(client, contractAddr)

	gotTokenType, err := con.TokenType(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mTokenType")
	}
	if gotTokenType != wantTokenType {
		return nil, nil, fmt.Errorf("got token type %s, want %s", gotTokenType.Hex(), wantTokenType.Hex())
	}

	gotAmount, err := con.Amount(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mAmount")
	}
	if gotAmount.Cmp(wantAmount) != 0 {
		return nil, nil, fmt.Errorf("got amount %s, want %s", gotAmount, wantAmount)
	}

	gotCollateral, err := con.Collateral(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mCollateral")
	}
	if gotCollateral.Cmp(wantCollateral) != 0 {
		return nil, nil, fmt.Errorf("got collateral %s, want %s", gotCollateral, wantCollateral)
	}

	gotClearRoot, err := con.ClearRoot(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mClearRoot")
	}
	if gotClearRoot != wantClearRoot {
		return nil, nil, fmt.Errorf("got clear root %x, want %x", gotClearRoot[:], wantClearRoot[:])
	}

	gotCipherRoot, err := con.CipherRoot(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mCipherRoot")
	}
	if gotCipherRoot != wantCipherRoot {
		return nil, nil, fmt.Errorf("got cipher root %x, want %x", gotCipherRoot[:], wantCipherRoot[:])
	}

	gotRevealDeadline, err := con.RevealDeadline(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mRevealDeadline")
	}
	if gotRevealDeadline.Unix() != wantRevealDeadline.Unix() { // lop off fractional seconds
		return nil, nil, fmt.Errorf("reveal deadline is %s, want %s", gotRevealDeadline, wantRevealDeadline)
	}

	if gotRevealDeadline.Sub(now) < minRevealDur {
		return nil, nil, fmt.Errorf("reveal deadline of %s is too soon, or in the past", gotRevealDeadline)
	}

	gotRefundDeadline, err := con.RefundDeadline(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting mRefundDeadline")
	}
	if gotRefundDeadline.Unix() != wantRefundDeadline.Unix() { // lop off fractional seconds from wantRefundDeadline
		return nil, nil, fmt.Errorf("refund deadline is %s, want %s", gotRefundDeadline, wantRefundDeadline)
	}

	paidAmount, err := con.Paid(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "checking paid amount")
	}
	if paidAmount.Cmp(wantAmount) < 0 {
		return nil, nil, fmt.Errorf("contract balance is %s, want %s", paidAmount, wantAmount)
	}

	if !IsETH(wantTokenType) {
		tokenInstance := erc20ABI.Instance(client, wantTokenType)
		if _, err := bind.Transact(tokenInstance, seller, erc20ABI.PackApprove(contractAddr, wantCollateral)); err != nil {
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

	revealTx, err := bind.Transact(con.BoundContract, revealTxOpts, treddABI.PackReveal(key))
	if err != nil {
		return nil, nil, errors.Wrap(err, "invoking ClaimPayment")
	}

	receipt, err := waitMined(ctx, client, revealTx)
	return con.BoundContract, receipt, errors.Wrap(err, "waiting for reveal tx to be mined")
}

// ClaimPayment constructs a seller-claims-payment transaction,
// rehydrating and invoking a Tredd contract from the utxo state (identified by the information in r).
func ClaimPayment(
	ctx context.Context,
	client clientType,
	seller *bind.TransactOpts,
	contractAddr common.Address,
) (*types.Receipt, error) {
	con := treddABI.Instance(client, contractAddr)
	tx, err := bind.Transact(con, seller, treddABI.PackClaimPayment())
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
	con *bind.BoundContract,
	index uint64,
	cipherChunk []byte,
	clearHash [32]byte,
	cipherProof, clearProof merkle.Proof,
) (*types.Receipt, error) {
	var (
		treddCipherProof = contract.Proof(cipherProof)
		treddClearProof  = contract.Proof(clearProof)
	)

	tx, err := bind.Transact(con, buyer, treddABI.PackRefund(index, cipherChunk, clearHash, treddCipherProof, treddClearProof))
	if err != nil {
		return nil, errors.Wrap(err, "invoking Refund")
	}
	return waitMined(ctx, client, tx)
}

func waitMined(ctx context.Context, client clientType, tx *types.Transaction) (*types.Receipt, error) {
	if c, ok := client.(committer); ok {
		c.Commit()
	}
	rcpt, err := bind.WaitMined(ctx, client, tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "waiting for transaction to be mined")
	}
	if rcpt.Status == types.ReceiptStatusFailed {
		return rcpt, fmt.Errorf("transaction %s failed", tx.Hash().Hex())
	}
	return rcpt, nil
}

var ethAddr common.Address

func IsETH(tokenType common.Address) bool {
	return tokenType == ethAddr
}
