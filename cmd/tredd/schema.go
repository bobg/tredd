package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const schema = `
CREATE TABLE IF NOT EXISTS transfers (
  transfer_id BLOB NOT NULL PRIMARY KEY,
  clear_root BLOB NOT NULL,
  cipher_root BLOB NOT NULL,
  contract_addr BLOB, -- NULL until discovered
  token_type BLOB NOT NULL,
  amount TEXT NOT NULL,
  collateral TEXT NOT NULL,
  buyer BLOB NOT NULL,
  reveal_deadline_secs INTEGER NOT NULL,
  refund_deadline_secs INTEGER NOT NULL,
  key BLOB NOT NULL
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
