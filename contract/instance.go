package contract

import (
	"context"
	"math/big"
	"time"

	"github.com/bobg/errors"
	"github.com/bobg/merkle/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

var treddABI = NewTredd()

func TokenType(ctx context.Context, con *bind.BoundContract) (common.Address, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMTokenType(), treddABI.UnpackMTokenType)
}

func Amount(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMAmount(), treddABI.UnpackMAmount)
}

func Collateral(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMCollateral(), treddABI.UnpackMCollateral)
}

func ClearRoot(ctx context.Context, con *bind.BoundContract) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMClearRoot(), treddABI.UnpackMClearRoot)
}

func CipherRoot(ctx context.Context, con *bind.BoundContract) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMCipherRoot(), treddABI.UnpackMCipherRoot)
}

func RevealDeadline(ctx context.Context, con *bind.BoundContract) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	revealDeadlineSecs, err := bind.Call(con, callOpts, treddABI.PackMRevealDeadline(), treddABI.UnpackMRevealDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRevealDeadline")
	}
	return time.Unix(int64(revealDeadlineSecs), 0), nil
}

func RefundDeadline(ctx context.Context, con *bind.BoundContract) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	refundDeadlineSecs, err := bind.Call(con, callOpts, treddABI.PackMRefundDeadline(), treddABI.UnpackMRefundDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRefundDeadline")
	}
	return time.Unix(int64(refundDeadlineSecs), 0), nil
}

func Paid(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackPaid(), treddABI.UnpackPaid)
}

func CheckProofWithPrefixedChunk(ctx context.Context, con *bind.BoundContract, chunkProof merkle.Proof, idx uint64, refChunk []byte, chunkRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedChunk(Proof(chunkProof), idx, refChunk, chunkRoot)
	)
	return bind.Call(con, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedChunk)
}

func CheckProofWithPrefixedHash(ctx context.Context, con *bind.BoundContract, hashProof merkle.Proof, idx uint64, refHash [32]byte, hashRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedHash(Proof(hashProof), idx, refHash, hashRoot)
	)
	return bind.Call(con, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedHash)
}

func Decrypt(con *bind.BoundContract, cipher []byte) ([]byte, error) {
	callOpts := &bind.CallOpts{}
	return bind.Call(con, callOpts, treddABI.PackDecrypt(cipher, 0), treddABI.UnpackDecrypt)
}
