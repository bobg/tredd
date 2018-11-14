package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/golang/protobuf/proto"
)

func submitter(url string) func(prog []byte, version, runlimit int64) error {
	return func(prog []byte, version, runlimit int64) error {
		rawTx := &bc.RawTx{
			Version:  version,
			Runlimit: runlimit,
			Program:  prog,
		}
		bits, err := proto.Marshal(rawTx)
		if err != nil {
			return errors.Wrap(err, "serializing tx")
		}
		resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(bits))
		if err != nil {
			return errors.Wrap(err, "submitting tx")
		}
		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("status code %d when submitting tx", resp.StatusCode)
		}
		return nil
	}
}
