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
	Reserve(context.Context, int64, bc.Hash, time.Time) (Reservation, error)
}

type Reservation interface {
	UTXOs() []UTXO
	Change() int64
	Cancel()
}
