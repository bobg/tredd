package main

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/bobg/sqlutil"
	"github.com/bobg/tredd"
	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
	"github.com/chain/txvm/protocol/txbuilder/txresult"
)

func processBlock(ctx context.Context, dbtx *sql.Tx, b *bc.Block, pubkey ed25519.PublicKey) error {
	for _, bctx := range b.Transactions {
		txr := txresult.New(bctx)
		for _, inp := range txr.Inputs {
			if inp.Value == nil {
				continue
			}
			_, err := dbtx.ExecContext(ctx, "DELETE FROM utxos WHERE output_id = $1", inp.OutputID)
			if err != nil {
				return errors.Wrapf(err, "deleting consumed utxo %x", inp.OutputID)
			}
		}
		for _, out := range txr.Outputs {
			if out.Value == nil {
				continue
			}
			if len(out.Pubkeys) != 1 {
				continue
			}
			if !bytes.Equal(out.Pubkeys[0], pubkey) {
				continue
			}
			_, err := dbtx.ExecContext(ctx, "INSERT INTO utxos (output_id, asset_id, amount, anchor) VALUES ($1, $2, $3, $4)", out.OutputID, out.Value.AssetID, out.Value.Amount, out.Value.Anchor)
			if err != nil {
				return errors.Wrapf(err, "inserting new utxo %x", out.OutputID)
			}
		}
	}
	return nil
}

type reservation struct {
	r      *reserver
	id     int64
	change int64
}

var errInsufficientFunds = errors.New("insufficient funds")

type reserver struct {
	db *sql.DB
}

func (r *reserver) Reserve(ctx context.Context, amount int64, assetID bc.Hash, now, exp time.Time) (tredd.Reservation, error) {
	dbtx, err := r.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "beginning db transaction")
	}
	defer dbtx.Rollback()

	expMS := bc.Millis(exp)
	dbres, err := dbtx.ExecContext(ctx, "INSERT INTO reservations (expiration_ms) VALUES ($1)", expMS)
	if err != nil {
		return nil, errors.Wrap(err, "inserting new reservation into db")
	}
	reservationID, err := dbres.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "querying insert id of new reservation")
	}

	res := &reservation{r: r, id: reservationID}

	const q = `
		SELECT u.amount, u.output_id FROM utxos u
		WHERE u.asset_id = $1
			AND NOT EXISTS
			(SELECT 1 FROM reservations r, reservation_utxos ru
			 WHERE r.expiration_ms > $2
			   AND r.reservation_id = ru.reservation_id
			   AND ru.output_id = u.output_id)
	`
	nowMS := bc.Millis(now)
	var outputIDs [][]byte
	err = sqlutil.ForQueryRows(ctx, dbtx, q, assetID, nowMS, func(utxoAmount int64, outputID []byte) {
		if amount <= 0 {
			return
		}
		outputIDs = append(outputIDs, outputID)
		amount -= utxoAmount
	})
	if err != nil {
		return nil, errors.Wrap(err, "querying db")
	}
	if amount > 0 {
		return nil, errors.Wrapf(errInsufficientFunds, "reserving %d of %x", amount, assetID.Bytes())
	}
	res.change = -amount
	for _, outputID := range outputIDs {
		_, err = dbtx.ExecContext(ctx, "INSERT INTO reservation_utxos (reservation_id, output_id) VALUES ($1, $2)", reservationID, outputID)
		if err != nil {
			return nil, errors.Wrapf(err, "adding utxo %x to reservation %d", outputID, reservationID)
		}
	}
	err = dbtx.Commit()
	return res, errors.Wrap(err, "committing db transaction")
}

func (r *reservation) UTXOs(ctx context.Context) ([]tredd.UTXO, error) {
	const q = `
		SELECT u.amount, u.asset_id, u.anchor
			FROM utxos u, reservation_utxos ru
			WHERE u.output_id = ru.output_id AND ru.reservation_id = $1
	`

	var utxos []tredd.UTXO
	err := sqlutil.ForQueryRows(ctx, r.r.db, q, r.id, func(amount int64, assetID bc.Hash, anchor []byte) {
		u := &utxo{
			amount:  amount,
			assetID: assetID,
			anchor:  anchor,
		}
		utxos = append(utxos, u)
	})
	return utxos, errors.Wrap(err, "querying db")
}

func (r *reservation) Change(context.Context) (int64, error) {
	return r.change, nil
}

func (r *reservation) Cancel(ctx context.Context) error {
	_, err := r.r.db.ExecContext(ctx, "DELETE FROM reservation WHERE id = $1", r.id)
	return errors.Wrap(err, "canceling reservation")
}

type utxo struct {
	amount  int64
	assetID bc.Hash
	anchor  []byte
}

func (u *utxo) Amount() int64    { return u.amount }
func (u *utxo) AssetID() bc.Hash { return u.assetID }
func (u *utxo) Anchor() []byte   { return u.anchor }
