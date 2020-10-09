package main

import (
	"flag"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// These functions are duplicated from github.com/bobg/ninex.

func addKeyfilePassphrase(fs *flag.FlagSet) (keyfile, passphrase *string) {
	keyfile = fs.String("keyfile", "", "key file")
	passphrase = fs.String("passphrase", "", "passphrase")
	return keyfile, passphrase
}

func handleKeyfilePassphrase(keyfile, passphrase string) (*bind.TransactOpts, error) {
	keyReader, err := os.Open(keyfile)
	if err != nil {
		return nil, err
	}
	defer keyReader.Close()

	return bind.NewTransactor(keyReader, passphrase)
}
