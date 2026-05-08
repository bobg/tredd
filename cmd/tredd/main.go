package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"log"
	"net/http"
	"os"
	"path"

	"github.com/bobg/errors"
	"github.com/bobg/merkle/v2"
	"github.com/bobg/subcmd/v2"

	"github.com/bobg/tredd"
	"github.com/bobg/tredd/contract"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	return subcmd.Run(ctx, maincmd{}, os.Args[1:])
}

type maincmd struct{}

func (maincmd) Subcmds() subcmd.Map {
	return subcmd.Commands(
		"add", add, "add files to a content tree", subcmd.Params(
			"-dir", subcmd.String, ".", "root of content tree",
			"-type", subcmd.String, "", "MIME content type (default: inferred)",
		),
		"decrypt", decrypt, "decrypt data from stdin using the given key", subcmd.Params(
			"-key", subcmd.String, "", "decryption key (hex)",
		),
		"get", get, "propose payment and get a file from a tredd server", subcmd.Params(
			"-hash", subcmd.String, "", "clear-chunk Merkle root hash of requested file",
			"-token", subcmd.String, "", "token type (ERC20 hex address) of proposed payment, or omit for ETH",
			"-amount", subcmd.String, "1", "amount of proposed payment",
			"-collateral", subcmd.String, "1", "amount of proposed collateral",
			"-reveal", subcmd.Duration, 15*time.Minute, "time until reveal deadline, in time.ParseDuration format",
			"-refund", subcmd.Duration, 30*time.Minute, "time from reveal deadline until refund deadline, in time.ParseDuration format",
			"-server", subcmd.String, "", "base URL of tredd server",
			"-ethurl", subcmd.String, "", "base URL of Ethereum server",
			"-dir", subcmd.String, "", "root dir for file transfers",
			"-seller", subcmd.String, "", "seller address (hex)",
			"-keyfile", subcmd.String, "", "path to Ethereum keyfile for payment",
			"-passphrase", subcmd.String, "", "passphrase for Ethereum keyfile",
		),
		"serve", serve, "start a tredd server", subcmd.Params(
			"-addr", subcmd.String, "localhost:20544", "server listen address",
			"-dir", subcmd.String, ".", "root of content tree",
			"-db", subcmd.String, "", "file containing server-state db",
			"-ethurl", subcmd.String, "", "base URL of Ethereum server",
			"-keyfile", subcmd.String, "", "path to Ethereum keyfile for payment",
			"-passphrase", subcmd.String, "", "passphrase for Ethereum keyfile",
		),
		"abi", abi, "print the Tredd contract ABI", nil,
	)
}

func add(_ context.Context, dir, contentType string, args []string) error {
	for _, file := range args {
		if err := addFile(file, dir, contentType); err != nil {
			log.Printf("WARNING: while processing %s: %s", file, err)
		}
	}
	return nil
}

func addFile(file, dir, contentType string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "opening %s", file)
	}
	defer f.Close()

	var (
		tree  = merkle.NewTree(sha256.New())
		chunk [tredd.ChunkSize]byte
	)

	for index := uint64(0); ; index++ {
		n, err := io.ReadFull(f, chunk[:])
		if errors.Is(err, io.EOF) {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			return errors.Wrapf(err, "reading %s", file)
		}
		if index == 0 && contentType == "" {
			contentType = http.DetectContentType(chunk[:n])
		}

		clearHash := sha256.Sum256(chunk[:n])
		tree.Add(tredd.Prefix(index, clearHash[:]))
	}

	var clearRoot [32]byte
	copy(clearRoot[:], tree.Root())

	p, destName := clearRootPath(dir, clearRoot)

	if err := os.MkdirAll(p, 0700); err != nil {
		return errors.Wrapf(err, "creating dir %s", p)
	}

	f.Close()

	if err := os.WriteFile(path.Join(p, "content-type"), []byte(contentType), 0600); err != nil {
		return errors.Wrapf(err, "storing content type: %s", err)
	}

	f, err = os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "reopening %s", file)
	}
	defer f.Close()

	dest, err := os.Create(path.Join(p, destName))
	if err != nil {
		return errors.Wrapf(err, "creating destination %s", destName)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, f); err != nil {
		return errors.Wrapf(err, "copying %s to %s", file, destName)
	}

	fmt.Printf("added %s (content type %s) as %x\n", file, contentType, clearRoot)

	return nil
}

func clearRootPath(root string, clearRoot [32]byte) (dir, filename string) {
	dir = path.Join(root, fmt.Sprintf("%x/%x", clearRoot[0:1], clearRoot[1:2]))
	return dir, hex.EncodeToString(clearRoot[:])
}

func decrypt(_ context.Context, keyHex string, _ []string) error {
	var key [32]byte
	if _, err := hex.Decode(key[:], []byte(keyHex)); err != nil {
		return errors.Wrap(err, "decoding key")
	}
	for index := uint64(0); ; index++ {
		var buf [tredd.ChunkSize]byte
		n, err := io.ReadFull(os.Stdin, buf[:])
		if errors.Is(err, io.EOF) {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			return errors.Wrap(err, "reading cipher chunk")
		}
		tredd.Crypt(key, buf[:n], index)
		os.Stdout.Write(buf[:n])
	}

	return nil
}

func abi(_ context.Context, _ []string) error {
	_, err := fmt.Println(contract.TreddMetaData.ABI)
	return errors.Wrap(err, "printing ABI")
}
