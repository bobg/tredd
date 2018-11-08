package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/bobg/merkle"
	"github.com/bobg/tedd"
	"github.com/chain/txvm/errors"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: tedd add [-dir DIR] FILE ...")
	}
	switch os.Args[1] {
	case "add":
		add(os.Args[2:])
	default:
		log.Fatalf("unknown subcommand %s", os.Args[1])
	}
}

func add(args []string) {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	dir := fs.String("dir", ".", "root of content tree")
	err := fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fs.Args() {
		err = addFile(file, *dir)
		if err != nil {
			log.Print("WARNING: while processing %s: %s", file, err)
		}
	}
}

func addFile(file, dir string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "opening %s", file)
	}
	defer f.Close()

	tree := merkle.NewTree(sha256.New())

	var buf [binary.MaxVarintLen64 + tedd.ChunkSize]byte
	for i := uint64(0); ; i++ {
		m := binary.PutUvarint(buf[:], i)
		n, err := io.ReadFull(f, buf[m:m+tedd.ChunkSize])
		if err == io.EOF {
			// "The error is EOF only if no bytes were read."
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return errors.Wrapf(err, "reading %s", file)
		}
		tree.Add(buf[:m+n])
	}
	clearHash := tree.Root()
	p := path.Join(dir, fmt.Sprintf("%x/%x", clearHash[0:1], clearHash[1:2]))
	err = os.MkdirAll(p, 0700)
	if err != nil {
		return errors.Wrapf(err, "creating dir %s", p)
	}

	f.Close()

	f, err = os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "reopening %s", file)
	}
	defer f.Close()

	destName := hex.EncodeToString(clearHash)
	dest, err := os.Create(path.Join(p, destName))
	if err != nil {
		return errors.Wrapf(err, "creating destination %s", destName)
	}
	defer dest.Close()

	_, err = io.Copy(dest, f)
	if err != nil {
		return errors.Wrapf(err, "copying %s to %s", file, destName)
	}

	return nil
}
