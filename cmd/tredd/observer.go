package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/chain/txvm/crypto/ed25519"
	"github.com/chain/txvm/errors"
	"github.com/chain/txvm/protocol/bc"
)

// Observer observes a blockchain,
// tracking its height,
// the timestamp of the latest block,
// and the utxos for a particular pubkey.
// It also maintains a queue of timers
// and a callback function called on every tx of every block.
type observer struct {
	db     *sql.DB
	pubkey ed25519.PublicKey
	r      *reserver
	url    string

	mu    sync.Mutex // protects cb and queue
	cb    func(*bc.Tx)
	queue []timer // ordered by time
}

type timer struct {
	t time.Time
	f func()
}

func newObserver(db *sql.DB, pubkey ed25519.PublicKey, url string) *observer {
	return &observer{
		db:     db,
		pubkey: pubkey,
		r:      &reserver{db: db},
		url:    url,
	}
}

// runs as a goroutine until its context is canceled
func (o *observer) run(ctx context.Context) {
	var client http.Client

	for {
		height, err := o.height(ctx)
		if err != nil {
			log.Fatalf("getting blockchain height: %s", err)
		}

		log.Printf("requesting block at height %d", height+1)

		getBlockURL := fmt.Sprintf("%s?height=%d", o.url, height+1)
		req, err := http.NewRequest("GET", getBlockURL, nil)
		if err != nil {
			log.Fatalf("constructing get-block request for %s: %s", getBlockURL, err)
		}
		req = req.WithContext(ctx)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error getting block at height %d from %s, will retry in ~5 seconds: %s", height+1, o.url, err)

			t := time.NewTimer(5 * time.Second) // TODO: add jitter
			select {
			case <-t.C:
				continue

			case <-ctx.Done():
				// canceled, exit
				t.Stop()
				log.Print("context canceled, blockchain observer exiting")
				return
			}
		}

		err = func() error {
			defer resp.Body.Close()

			bits, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, "reading block")
			}
			b := new(bc.Block)
			err = b.FromBytes(bits)
			if err != nil {
				return errors.Wrap(err, "parsing block")
			}

			dbtx, err := o.db.Begin()
			if err != nil {
				return errors.Wrap(err, "beginning db transaction")
			}
			defer dbtx.Rollback()

			err = processBlock(ctx, dbtx, b, o.pubkey)
			if err != nil {
				return errors.Wrap(err, "updating reserver")
			}

			_, err = dbtx.ExecContext(ctx, "INSERT OR REPLACE INTO latest_block (singleton, height, timestamp_ms) VALUES (0, $1, $2)", b.Height, b.TimestampMs)
			if err != nil {
				return errors.Wrap(err, "storing block info")
			}

			err = dbtx.Commit()
			if err != nil {
				return errors.Wrap(err, "committing db transaction")
			}

			now := bc.FromMillis(b.TimestampMs)
			log.Printf("block time %s", now)

			o.mu.Lock()
			for len(o.queue) > 0 && !o.queue[0].t.After(now) {
				go o.queue[0].f()
				o.queue = o.queue[1:]
			}

			if o.cb != nil {
				for _, tx := range b.Transactions {
					tx := tx // Go loop-var pitfall
					go o.cb(tx)
				}
			}
			o.mu.Unlock()
			return nil
		}()
		if err != nil {
			log.Fatalf("processing block %d: %s", height+1, err)
		}
	}
}

func (o *observer) setcb(cb func(*bc.Tx)) {
	o.mu.Lock()
	o.cb = cb
	o.mu.Unlock()
}

func (o *observer) height(ctx context.Context) (uint64, error) {
	var height uint64
	err := o.db.QueryRowContext(ctx, "SELECT height FROM latest_block WHERE singleton = 0").Scan(&height)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return height, err
}

func (o *observer) now(ctx context.Context) (time.Time, error) {
	var timestampMS uint64
	err := o.db.QueryRowContext(ctx, "SELECT timestamp_ms FROM latest_block WHERE singleton = 0").Scan(&timestampMS)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "querying timestamp from db")
	}
	return bc.FromMillis(timestampMS), nil
}

func (o *observer) enqueue(t time.Time, f func()) {
	o.mu.Lock()
	o.queue = append(o.queue, timer{t: t, f: f})
	sort.Slice(o.queue, func(i, j int) bool {
		return o.queue[i].t.Before(o.queue[j].t)
	})
	o.mu.Unlock()
}
