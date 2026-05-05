package main

import (
	"context"
	"math/big"
	"os"

	"github.com/bobg/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// These functions are duplicated from github.com/bobg/ninex.

func handleKeyfilePassphrase(ctx context.Context, keyfile, passphrase string, chainID *big.Int) (*bind.TransactOpts, error) {
	keyJSON, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, errors.Wrapf(err, "reading keyfile %q", keyfile)
	}

	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return nil, errors.Wrapf(err, "decrypting keyfile %q with given passphrase", keyfile)
	}

	txOpts, err := bind.NewKeyedTransactor(key.PrivateKey, chainID), nil
	if err != nil {
		return nil, errors.Wrapf(err, "creating transactor from keyfile %q", keyfile)
	}

	txOpts.Context = ctx

	return txOpts, nil
}
