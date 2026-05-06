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

// Instance is a wrapper for [bind.BoundContract] that provides typed methods for calling some functions of the Tredd contract.
type Instance struct {
	*bind.BoundContract
}

// NewInstance creates a new bound instance of the Tredd contract at the given address.
func NewInstance(backend bind.ContractBackend, addr common.Address) Instance {
	boundContract := treddABI.Instance(backend, addr)
	return Instance{BoundContract: boundContract}
}

// TokenType returns the token type of the payment, which is either an ERC20 token address or the zero address for ETH.
func (inst Instance) TokenType(ctx context.Context) (common.Address, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackMTokenType(), treddABI.UnpackMTokenType)
}

// Amount returns the amount of the proposed payment.
func (inst Instance) Amount(ctx context.Context) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackMAmount(), treddABI.UnpackMAmount)
}

// Collateral returns the collateral amount requested by the buyer.
func (inst Instance) Collateral(ctx context.Context) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackMCollateral(), treddABI.UnpackMCollateral)
}

// ClearRoot returns the Merkle root hash of the cleartext chunks of the content.
func (inst Instance) ClearRoot(ctx context.Context) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackMClearRoot(), treddABI.UnpackMClearRoot)
}

// CipherRoot returns the Merkle root hash of the ciphertext chunks of the content.
func (inst Instance) CipherRoot(ctx context.Context) ([32]byte, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackMCipherRoot(), treddABI.UnpackMCipherRoot)
}

// RevealDeadline returns the reveal deadline.
func (inst Instance) RevealDeadline(ctx context.Context) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	revealDeadlineSecs, err := bind.Call(inst.BoundContract, callOpts, treddABI.PackMRevealDeadline(), treddABI.UnpackMRevealDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRevealDeadline")
	}
	return time.Unix(int64(revealDeadlineSecs), 0), nil
}

// RefundDeadline returns the refund deadline.
func (inst Instance) RefundDeadline(ctx context.Context) (time.Time, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	refundDeadlineSecs, err := bind.Call(inst.BoundContract, callOpts, treddABI.PackMRefundDeadline(), treddABI.UnpackMRefundDeadline)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "getting mRefundDeadline")
	}
	return time.Unix(int64(refundDeadlineSecs), 0), nil
}

// Paid returns the amount paid by the buyer so far.
func (inst Instance) Paid(ctx context.Context) (*big.Int, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackPaid(), treddABI.UnpackPaid)
}

// CheckProofWithPrefixedChunk checks a Merkle proof of a specific cleartext chunk against the clear root.
func (inst Instance) CheckProofWithPrefixedChunk(ctx context.Context, chunkProof merkle.Proof, idx uint64, refChunk []byte, chunkRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedChunk(Proof(chunkProof), idx, refChunk, chunkRoot)
	)
	return bind.Call(inst.BoundContract, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedChunk)
}

// CheckProofWithPrefixedHash checks a Merkle proof of a specific ciphertext chunk against the cipher root.
// This is used for proofs of ciphertext chunks, since the contract only knows the ciphertext hashes.
func (inst Instance) CheckProofWithPrefixedHash(ctx context.Context, hashProof merkle.Proof, idx uint64, refHash [32]byte, hashRoot [32]byte) (bool, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		callData = treddABI.PackCheckProofWithPrefixedHash(Proof(hashProof), idx, refHash, hashRoot)
	)
	return bind.Call(inst.BoundContract, callOpts, callData, treddABI.UnpackCheckProofWithPrefixedHash)
}

// Decrypt decrypts a cipher chunk.
func (inst Instance) Decrypt(cipher []byte) ([]byte, error) {
	callOpts := &bind.CallOpts{}
	return bind.Call(inst.BoundContract, callOpts, treddABI.PackDecrypt(cipher, 0), treddABI.UnpackDecrypt)
}
