package main

import (
	"context"
	"database/sql"

	"github.com/chain/txvm/errors"
	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS utxos (
  output_id BLOB NOT NULL PRIMARY KEY,
  asset_id BLOB NOT NULL,
  amount INTEGER NOT NULL,
  anchor BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS utxos_asset_id ON utxos (asset_id);
CREATE UNIQUE INDEX IF NOT EXISTS utxos_anchor ON utxos (anchor);

CREATE TABLE IF NOT EXISTS reservations (
  reservation_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  expiration_ms INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS reservation_utxos (
  reservation_id INTEGER NOT NULL REFERENCES reservations ON DELETE CASCADE,
  output_id BLOB NOT NULL REFERENCES utxos
);

CREATE TABLE IF NOT EXISTS latest_block (
  singleton INTEGER NOT NULL PRIMARY KEY DEFAULT 0 CHECK (singleton = 0),
  height INTEGER NOT NULL,
  timestamp_ms INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS transfer_records (
  transfer_id BLOB NOT NULL PRIMARY KEY,
  reveal_deadline_ms INTEGER NOT NULL,
  refund_deadline_ms INTEGER NOT NULL,
  cipher_root BLOB NOT NULL,
  clear_root BLOB NOT NULL,
  amount INTEGER NOT NULL,
  asset_id BLOB NOT NULL,
  anchor1 BLOB,
  anchor2 BLOB,
  key BLOB,
  output_id BLOB,
  seller BLOB NOT NULL,
  buyer BLOB
);
`

func setSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}

func openDB(ctx context.Context, filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s", filename)
	}
	err = setSchema(ctx, db)
	return db, errors.Wrap(err, "setting db schema")
}
