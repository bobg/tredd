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

// TokenType returns the token type of the payment, which is either an ERC20 token address or the zero address for ETH.
func TokenType(ctx context.Context, con *bind.BoundContract) (common.Address, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMTokenType(), treddABI.UnpackMTokenType)
}

// Amount returns the amount of the proposed payment.
func Amount(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMAmount(), treddABI.UnpackMAmount)
}

// Collateral returns the collateral amount requested by the buyer.
func Collateral(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMCollateral(), treddABI.UnpackMCollateral)
}

// ClearRoot returns the Merkle root hash of the cleartext chunks of the content.
func ClearRoot(ctx context.Context, con *bind.BoundContract) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMClearRoot(), treddABI.UnpackMClearRoot)
}

// CipherRoot returns the Merkle root hash of the ciphertext chunks of the content.
func CipherRoot(ctx context.Context, con *bind.BoundContract) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackMCipherRoot(), treddABI.UnpackMCipherRoot)
}

// RevealDeadline returns the reveal deadline.
func RevealDeadline(ctx context.Context, con *bind.BoundContract) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	revealDeadlineSecs, err := bind.Call(con, callOpts, treddABI.PackMRevealDeadline(), treddABI.UnpackMRevealDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRevealDeadline")
	}
	return time.Unix(int64(revealDeadlineSecs), 0), nil
}

// RefundDeadline returns the refund deadline.
func RefundDeadline(ctx context.Context, con *bind.BoundContract) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	refundDeadlineSecs, err := bind.Call(con, callOpts, treddABI.PackMRefundDeadline(), treddABI.UnpackMRefundDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRefundDeadline")
	}
	return time.Unix(int64(refundDeadlineSecs), 0), nil
}

// Paid returns the amount paid by the buyer so far.
func Paid(ctx context.Context, con *bind.BoundContract) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(con, callOpts, treddABI.PackPaid(), treddABI.UnpackPaid)
}

// CheckProofWithPrefixedChunk checks a Merkle proof of a specific cleartext chunk against the clear root.
func CheckProofWithPrefixedChunk(ctx context.Context, con *bind.BoundContract, chunkProof merkle.Proof, idx uint64, refChunk []byte, chunkRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedChunk(Proof(chunkProof), idx, refChunk, chunkRoot)
	)
	return bind.Call(con, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedChunk)
}

// CheckProofWithPrefixedHash checks a Merkle proof of a specific ciphertext chunk against the cipher root.
// This is used for proofs of ciphertext chunks, since the contract only knows the ciphertext hashes.
func CheckProofWithPrefixedHash(ctx context.Context, con *bind.BoundContract, hashProof merkle.Proof, idx uint64, refHash [32]byte, hashRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedHash(Proof(hashProof), idx, refHash, hashRoot)
	)
	return bind.Call(con, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedHash)
}

// Decrypt decrypts a cipher chunk.
func Decrypt(con *bind.BoundContract, cipher []byte) ([]byte, error) {
	callOpts := &bind.CallOpts{}
	return bind.Call(con, callOpts, treddABI.PackDecrypt(cipher, 0), treddABI.UnpackDecrypt)
}
