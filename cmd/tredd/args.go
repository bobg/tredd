package main

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// These functions are duplicated from github.com/bobg/ninex.

func handleKeyfilePassphrase(keyfile, passphrase string) (*bind.TransactOpts, error) {
	keyReader, err := os.Open(keyfile)
	if err != nil {
		return nil, err
	}
	defer keyReader.Close()

	return bind.NewTransactor(keyReader, passphrase)
}
