package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const schema = `
CREATE TABLE IF NOT EXISTS transfer_records (
  transfer_id BLOB NOT NULL PRIMARY KEY,
  reveal_deadline_ms INTEGER NOT NULL,
  refund_deadline_ms INTEGER NOT NULL,
  cipher_root BLOB NOT NULL,
  clear_root BLOB NOT NULL,
  amount INTEGER NOT NULL,
  token_type BLOB NOT NULL,
  key BLOB,
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
