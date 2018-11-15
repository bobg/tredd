package main

import (
	"context"
	"encoding/binary"
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
	"github.com/coreos/bbolt"
)

// Observer observes a blockchain,
// tracking its height,
// the timestamp of the latest block,
// and the utxos for a particular pubkey.
// It also maintains a queue of timers
// and a callback function called on every tx of every block.
type observer struct {
	db     *bbolt.DB
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

func newObserver(db *bbolt.DB, pubkey ed25519.PublicKey, url string) *observer {
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
		height, err := o.height()
		if err != nil {
			log.Fatalf("getting blockchain height: %s", err)
		}

		getBlockURL := fmt.Sprintf("%s?height=%d", o.url, height+1)
		req, err := http.NewRequest("GET", getBlockURL, nil)
		if err != nil {
			log.Fatalf("constructing get-block request for %s: %s", getBlockURL, err)
		}
		req = req.WithContext(ctx)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("error getting block at height %d from %s, will retry in ~5 seconds: %s", height+1, o.url, err)

			t := time.NewTimer(5 * time.Second) // xxx add jitter
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

			err = o.db.Update(func(tx *bbolt.Tx) error {
				root, err := tx.CreateBucketIfNotExists([]byte("root"))
				if err != nil {
					return errors.Wrap(err, "getting/creating root bucket")
				}

				var heightBuf [binary.MaxVarintLen64]byte
				m := binary.PutUvarint(heightBuf[:], b.Height)
				err = root.Put([]byte("height"), heightBuf[:m])
				if err != nil {
					return errors.Wrap(err, "storing blockchain height")
				}

				var nowMSBuf [binary.MaxVarintLen64]byte
				m = binary.PutUvarint(nowMSBuf[:], b.TimestampMs)
				err = root.Put([]byte("nowMS"), nowMSBuf[:m])
				if err != nil {
					return errors.Wrap(err, "storing blockchain time")
				}

				return nil
			})
			if err != nil {
				return errors.Wrap(err, "storing block info")
			}

			err = o.db.Update(func(tx *bbolt.Tx) error {
				return processBlock(tx, b, o.pubkey)
			})
			if err != nil {
				return errors.Wrap(err, "updating reserver")
			}

			now := bc.FromMillis(b.TimestampMs)

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

func (o *observer) height() (uint64, error) {
	var height uint64
	err := o.db.View(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return errors.Wrap(err, "getting/creating root bucket")
		}
		heightBits := root.Get([]byte("height"))
		if len(heightBits) == 0 {
			return nil
		}
		var n int
		height, n = binary.Uvarint(heightBits)
		if n < 1 {
			return fmt.Errorf("parsing blockchain height")
		}
		return nil
	})
	return height, errors.Wrap(err, "getting blockchain height")
}

func (o *observer) now() (time.Time, error) {
	var result time.Time
	err := o.db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte("root"))      // xxx check
		nowMSBits := root.Get([]byte("nowMS")) // xxx check
		nowMS, n := binary.Uvarint(nowMSBits)
		if n < 1 {
			return fmt.Errorf("parsing blockchain time")
		}
		result = bc.FromMillis(nowMS)
		return nil
	})
	return result, errors.Wrap(err, "getting blockchain time")
}

func (o *observer) enqueue(t time.Time, f func()) {
	o.mu.Lock()
	o.queue = append(o.queue, timer{t: t, f: f})
	sort.Slice(o.queue, func(i, j int) bool {
		return o.queue[i].t.Before(o.queue[j].t)
	})
	o.mu.Unlock()
}
