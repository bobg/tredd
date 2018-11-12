package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/bobg/tedd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txbuilder/txresult"
	"github.com/coreos/bbolt"
)

type reserver struct {
	pubkey ed25519.PublicKey
	db     *bbolt.DB
}

func (r *reserver) processBlock(b *bc.Block) error {
	return r.db.Update(func(dbtx *bbolt.Tx) error {
		utxos, err := dbtx.CreateBucketIfNotExists([]byte("utxos"))
		if err != nil {
			return errors.Wrapf(err, "getting/creating utxos db bucket")
		}
		for _, bctx := range b.Transactions {
			txr := txresult.New(bctx)
			for _, inp := range txr.Inputs {
				if inp.Value == nil {
					continue
				}
				asset := utxos.Bucket(inp.Value.AssetID.Bytes())
				if asset == nil {
					return errors.Wrapf(err, "asset bucket %x not found", inp.Value.AssetID.Bytes())
				}
				err = asset.DeleteBucket(inp.OutputID.Bytes())
				if err != nil {
					return errors.Wrapf(err, "deleting bucket for utxo %x", inp.OutputID.Bytes())
				}
			}
			for _, out := range txr.Outputs {
				if out.Value == nil {
					continue
				}
				if len(out.Pubkeys) != 1 {
					continue
				}
				if !bytes.Equal(out.Pubkeys[0], r.pubkey) {
					continue
				}
				asset, err := utxos.CreateBucketIfNotExists(out.Value.AssetID.Bytes())
				if err != nil {
					return errors.Wrapf(err, "creating asset ID bucket %x", out.Value.AssetID.Bytes())
				}
				utxoBucket, err := asset.CreateBucket(out.OutputID.Bytes())
				if err != nil {
					return errors.Wrapf(err, "creating utxo bucket %x", out.OutputID.Bytes())
				}
				err = utxoBucket.Put([]byte("anchor"), out.Value.Anchor)
				if err != nil {
					return errors.Wrapf(err, "storing anchor for utxo %x", out.OutputID.Bytes())
				}
				var amountbuf [binary.MaxVarintLen64]byte
				m := binary.PutVarint(amountbuf[:], int64(out.Value.Amount)) // xxx range checking
				err = utxoBucket.Put([]byte("amount"), amountbuf[:m])
				if err != nil {
					return errors.Wrapf(err, "storing amount for utxo %x", out.OutputID.Bytes())
				}
			}
		}
		return nil
	})
}

type reservation struct {
	r         *reserver
	utxos     []tedd.UTXO // 1:1 with outputIDs
	outputIDs [][]byte    // 1:1 with utxos
	change    int64
}

var errInsufficientFunds = errors.New("insufficient funds")

func (r *reserver) Reserve(_ context.Context, amount int64, assetID bc.Hash, now, exp time.Time) (tedd.Reservation, error) {
	res := &reservation{r: r}
	err := r.db.Update(func(tx *bbolt.Tx) error {
		utxos := tx.Bucket([]byte("utxos"))
		if utxos == nil {
			return errInsufficientFunds
		}
		asset := utxos.Bucket(assetID.Bytes())
		if asset == nil {
			return errInsufficientFunds
		}
		c := asset.Cursor()
		for outputID, _ := c.First(); amount > 0 && outputID != nil; outputID, _ = c.Next() {
			utxoBucket := asset.Bucket(outputID)
			if utxoBucket == nil {
				return errInsufficientFunds
			}
			utxoExpBytes := utxoBucket.Get([]byte("expiration"))
			if len(utxoExpBytes) > 0 {
				var utxoExp time.Time
				err := utxoExp.UnmarshalBinary(utxoExpBytes)
				if err != nil {
					return errors.Wrapf(err, "parsing expiration time of reserved utxo %x", outputID)
				}
				if now.Before(utxoExp) {
					// utxo is reserved
					continue
				}
			}
			utxoAmount, n := binary.Varint(utxoBucket.Get([]byte("amount")))
			if n < 1 {
				return fmt.Errorf("cannot parse amount in utxo %x", outputID)
			}
			utxoAnchor := utxoBucket.Get([]byte("anchor"))
			u := &utxo{
				amount:  utxoAmount,
				assetID: assetID,
				anchor:  utxoAnchor,
			}
			res.utxos = append(res.utxos, u)
			res.outputIDs = append(res.outputIDs, outputID)
			amount -= utxoAmount
		}
		res.change = -amount
		for _, o := range res.outputIDs {
			expBytes, err := exp.MarshalBinary()
			if err != nil {
				return errors.Wrapf(err, "storing reservation expiration time in utxo %x", o)
			}
			utxoBucket := asset.Bucket(o)
			utxoBucket.Put([]byte("expiration"), expBytes)
		}
		return nil
	})
	return res, err
}

func (r *reservation) UTXOs() []tedd.UTXO {
	return r.utxos
}

func (r *reservation) Change() int64 {
	return r.change
}

func (r *reservation) Cancel(context.Context) error {
	return r.r.db.Update(func(tx *bbolt.Tx) error {
		utxos := tx.Bucket([]byte("utxos"))
		if utxos == nil {
			return nil
		}
		asset := utxos.Bucket(r.utxos[0].AssetID().Bytes())
		if asset == nil {
			return nil
		}
		for _, o := range r.outputIDs {
			utxoBucket := asset.Bucket(o)
			if utxoBucket == nil {
				continue
			}
			utxoBucket.Delete([]byte("expiration"))
		}
		return nil
	})
}

type utxo struct {
	amount  int64
	assetID bc.Hash
	anchor  []byte
}

func (u *utxo) Amount() int64    { return u.amount }
func (u *utxo) AssetID() bc.Hash { return u.assetID }
func (u *utxo) Anchor() []byte   { return u.anchor }
