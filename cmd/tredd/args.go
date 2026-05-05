package main

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// These functions are duplicated from github.com/bobg/ninex.

func handleKeyfilePassphrase(keyfile, passphrase string, chainID *big.Int) (*bind.TransactOpts, error) {
	keyJSON, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactor(key.PrivateKey, chainID), nil
}
