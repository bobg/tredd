package tedd

import (
	"context"
	"time"

	"github.com/chain/txvm/protocol/bc"
)

type UTXO interface {
	Amount() int64
	AssetID() bc.Hash
	Anchor() []byte
}

type Reserver interface {
	Reserve(ctx context.Context, amount int64, assetID bc.Hash, now, exp time.Time) (Reservation, error)
}

type Reservation interface {
	UTXOs() []UTXO
	Change() int64
	Cancel(context.Context) error
}
