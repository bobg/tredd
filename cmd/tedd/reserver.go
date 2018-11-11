package main

import (
	"context"
	"encoding/binary"
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

func (r *reserver) add(outputID bc.Hash, value txresult.Value) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		utxos, err := tx.CreateBucketIfNotExists([]byte("utxos"))
		if err != nil {
			return errors.Wrap(err, "creating utxos bucket")
		}
		asset, err := utxos.CreateBucketIfNotExists(value.AssetID.Bytes())
		if err != nil {
			return errors.Wrapf(err, "creating asset ID bucket %x", value.AssetID.Bytes())
		}
		utxoBucket, err := asset.CreateBucket(outputID.Bytes())
		if err != nil {
			return errors.Wrapf(err, "creating utxo bucket %x", outputID.Bytes())
		}
		err = utxoBucket.Put([]byte("anchor"), value.Anchor)
		if err != nil {
			return errors.Wrapf(err, "storing anchor for utxo %x", outputID.Bytes())
		}
		var amountbuf [binary.MaxVarintLen64]byte
		m := binary.PutUvarint(amountbuf[:], value.Amount)
		err = utxoBucket.Put([]byte("amount"), amountbuf[:m])
		if err != nil {
			return errors.Wrapf(err, "storing amount for utxo %x", outputID.Bytes())
		}
		return nil
	})
}

func (r *reserver) remove(outputID, assetID bc.Hash) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		utxos := tx.Bucket([]byte("utxos"))
		if utxos == nil {
			// xxx err
		}
		asset := utxos.Bucket(assetID.Bytes())
		if asset == nil {
			// xxx err
		}
		return asset.DeleteBucket(outputID.Bytes())
	})
}

type reservation struct {
	utxos  []tedd.UTXO
	change int64
}

func (r *reserver) Reserve(_ context.Context, amount int64, assetID bc.Hash, now, exp time.Time) (tedd.Reservation, error) {
	var res reservation
	err := r.db.Update(func(tx *bbolt.Tx) error {
		utxos := tx.Bucket([]byte("utxos"))
		if utxos == nil {
			// xxx err
		}
		asset := utxos.Bucket(assetID.Bytes())
		if asset == nil {
			// xxx err
		}
		c := asset.Cursor()
		for amount > 0 {
			outputID, _ := c.Next() // xxx do we have to do c.First() first?
			utxoBucket := asset.Bucket(outputID)
			if utxoBucket == nil {
				// xxx err
			}
			utxoExp := utxoBucket.Get([]byte("expiration"))
			if len(utxoExp) > 0 {
				// xxx parse utxoExp into a time.Time
				// xxx compare with now, skip if utxo is reserved
			}
			utxoAmount, n := binary.Uvarint(utxoBucket.Get([]byte("amount")))
			if n < 1 {
				// xxx err
			}
			utxoAnchor := utxoBucket.Get([]byte("anchor"))
			if len(anchor) != 32 {
				// xxx err
			}
			u := &utxo{
				amount:  utxoAmount,
				assetID: assetID,
				anchor:  utxoAnchor,
			}
			res.utxos = append(res.utxos, u)
			amount -= utxoAmount
		}
		res.change = -amount
		for _, u := range res.utxos {
			// xxx set expiration of this utxo
		}
		return nil
	})
	return &res, err
}
