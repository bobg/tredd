package tredd

import (
	"context"
	"time"

	"github.com/chain/txvm/protocol/bc"
)

// UTXO is the type of an unspent output in the blockchain's UTXO set.
type UTXO interface {
	Amount() int64
	AssetID() bc.Hash
	Anchor() []byte
}

// Reserver can reserve UTXOs for spending before a given expiration time.
// A UTXO, once reserved, will not appear in another Reservation until/unless the first reservation expires or is canceled.
type Reserver interface {
	Reserve(ctx context.Context, amount int64, assetID bc.Hash, now, exp time.Time) (Reservation, error)
}

// Reservation is the result of reserving some UTXOs with a Reserver.
type Reservation interface {
	UTXOs() []UTXO
	Change() int64
	Cancel(context.Context) error
}
