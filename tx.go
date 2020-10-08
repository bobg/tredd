package tredd

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/bobg/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

// ProposePayment publishes a new instance of the Tredd contract instantiated with the given parameters.
func ProposePayment(
	ctx context.Context,
	client *ethclient.Client, // see ethclient.Dial
	buyer *bind.TransactOpts, // see bind.NewTransactor
	seller common.Address,
	amount int64,
	tokenType string, // TODO: how to specify the token type?
	clearRoot, cipherRoot [32]byte,
	revealDeadline, refundDeadline time.Time,
) (*types.Receipt, error) {
	parsed, err := abi.JSON(strings.NewReader(TreddABI))
	if err != nil {
		return nil, errors.Wrap(err, "parsing contract JSON to ABI")
	}

	_, tx, _, err := bind.DeployContract(buyer, parsed, common.FromHex(TreddBin), client)
	if err != nil {
		return nil, errors.Wrap(err, "deploying contract")
	}

	// Wait for tx to be mined on-chain.
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, errors.Wrap(err, "awaiting contract-deployment receipt")
	}

	// TODO: store contractAddr

	return receipt, nil
}

// RevealKey updates a Tredd contract on-chain by adding the decryption key.
// TODO: Must also supply collateral.
func RevealKey(
	ctx context.Context,
	client *ethclient.Client, // see ethclient.Dial
	seller *bind.TransactOpts, // see bind.NewTransactor
	contractAddr common.Address,
	key [32]byte,
	wantClearRoot, wantCipherRoot [32]byte,
	wantRevealDeadline, wantRefundDeadline time.Time,
) (*types.Receipt, error) {
	// TODO: read values from the on-chain contract, verify they match the "want" parameters
	con, err := NewTredd(contractAddr, client)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating deployed contract")
	}
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

	var (
		bigIndex         = big.NewInt(index)
		treddCipherProof = toTreddProof(cipherProof)
		treddClearProof  = toTreddProof(clearProof)
	)

	tx, err := con.Refund(buyer, bigIndex, cipherChunk, clearHash, treddCipherProof, treddClearProof)
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
