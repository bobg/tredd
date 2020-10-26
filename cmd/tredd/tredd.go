package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/bobg/merkle"
	"github.com/pkg/errors"

	"github.com/bobg/tredd"
	"github.com/bobg/tredd/contract"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: tredd add [-dir DIR] FILE ...")
	}
	switch os.Args[1] {
	case "add":
		add(os.Args[2:])
	case "decrypt":
		decrypt(os.Args[2:])
	case "get":
		get(os.Args[2:])
	case "serve":
		serve(os.Args[2:])
	case "abi":
		fmt.Println(contract.TreddABI)
	default:
		log.Fatalf("unknown subcommand %s", os.Args[1])
	}
}

func add(args []string) {
	fs := flag.NewFlagSet("", flag.PanicOnError)

	var (
		dir         = fs.String("dir", ".", "root of content tree")
		contentType = fs.String("type", "", "MIME content type (default: inferred)")
	)
	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fs.Args() {
		err = addFile(file, *dir, *contentType)
		if err != nil {
			log.Printf("WARNING: while processing %s: %s", file, err)
		}
	}
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
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return errors.Wrapf(err, "reading %s", file)
		}
		if index == 0 && contentType == "" {
			contentType = http.DetectContentType(chunk[:n])
		}

		prefixedClearChunk := tredd.PrefixChunk(index, chunk[:n])
		tree.Add(prefixedClearChunk)
	}

	var clearHash [32]byte
	copy(clearHash[:], tree.Root())

	p, destName := clearHashPath(dir, clearHash)

	err = os.MkdirAll(p, 0700)
	if err != nil {
		return errors.Wrapf(err, "creating dir %s", p)
	}

	f.Close()

	err = ioutil.WriteFile(path.Join(p, "content-type"), []byte(contentType), 0600)
	if err != nil {
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

	_, err = io.Copy(dest, f)
	if err != nil {
		return errors.Wrapf(err, "copying %s to %s", file, destName)
	}

	fmt.Printf("added %s (content type %s) as %x\n", file, contentType, clearHash)

	return nil
}

func clearHashPath(root string, clearHash [32]byte) (dir, filename string) {
	dir = path.Join(root, fmt.Sprintf("%x/%x", clearHash[0:1], clearHash[1:2]))
	return dir, hex.EncodeToString(clearHash[:])
}

func decrypt(args []string) {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	keyHex := fs.String("key", "", "decryption key (hex)")
	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	var key [32]byte
	_, err = hex.Decode(key[:], []byte(*keyHex))
	if err != nil {
		log.Fatal(err)
	}
	for index := uint64(0); ; index++ {
		var buf [tredd.ChunkSize]byte
		n, err := io.ReadFull(os.Stdin, buf[:])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			log.Fatal(err)
		}
		tredd.Crypt(key, buf[:n], index)
		os.Stdout.Write(buf[:n])
	}
}
