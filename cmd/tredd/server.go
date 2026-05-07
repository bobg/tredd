package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/bobg/errors"
	"github.com/bobg/mid"
	"github.com/bobg/seqs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bobg/tredd/contract"

	"github.com/bobg/tredd"
)

var serverTreddABI = contract.NewTredd()

func serve(
	ctx context.Context,
	addr string,
	dir string,
	dbFile string,
	ethURL string,
	keyfile string,
	passphrase string,
	_ []string,
) error {
	db, err := openDB(ctx, dbFile)
	if err != nil {
		return errors.Wrapf(err, "opening db file %s", dbFile)
	}
	defer db.Close()

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return errors.Wrapf(err, "dialing Ethereum service at %s", ethURL)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return errors.Wrap(err, "getting chain ID")
	}

	seller, err := handleKeyfilePassphrase(ctx, keyfile, passphrase, chainID)
	if err != nil {
		return errors.Wrap(err, "handling keyfile and passphrase")
	}

	s := &server{
		db:     db,
		dir:    dir,
		seller: seller,
		client: client,
	}

	rows, errptr := seqs.SQL[[]byte](ctx, db, "SELECT transfer_id FROM transfer_records")
	for transferID := range rows {
		log.Printf("queueing claim-payment callback for transfer %x", transferID)
		if err := s.queueClaimPayment(ctx, transferID); err != nil {
			return errors.Wrapf(err, "queueing claim-payment callback for transfer ID %x", transferID)
		}
	}
	if err := *errptr; err != nil {
		return errors.Wrap(err, "getting transfer records")
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "listening on %s", addr)
	}

	log.Printf("listening on %s", listener.Addr())

	mux := http.NewServeMux()
	mux.Handle("/request", mid.Err(s.serve))
	mux.Handle("/propose-payment", mid.Err(s.revealKey))

	h := http.Server{
		Handler: mux,
	}

	err = h.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return errors.Wrap(err, "serving requests")
}

type server struct {
	db     *sql.DB // transfer records and blockchain info
	dir    string  // content
	seller *bind.TransactOpts
	client *ethclient.Client
}

type serverRecord struct {
	transferID                     [32]byte
	contractAddr                   *common.Address // nil until discovered
	tokenType                      common.Address
	amount, collateral             *big.Int
	revealDeadline, refundDeadline time.Time
	buyer                          common.Address
	key, clearRoot, cipherRoot     [32]byte
}

const (
	minRevealDur = 10 * time.Minute
	maxRefundDur = time.Hour
)

var big0 = big.NewInt(0)

func (s *server) serve(w http.ResponseWriter, req *http.Request) error {
	var (
		buyerHex              = req.FormValue("buyer")
		clearRootHex          = req.FormValue("clearroot")
		tokenTypeHex          = req.FormValue("token")
		amountStr             = req.FormValue("amount")
		collateralStr         = req.FormValue("collateral")
		revealDeadlineSecsStr = req.FormValue("revealdeadline")
		refundDeadlineSecsStr = req.FormValue("refunddeadline")
	)

	buyer := common.HexToAddress(buyerHex)

	var clearRoot [32]byte
	if _, err := hex.Decode(clearRoot[:], []byte(clearRootHex)); err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "decoding clear root")}
	}

	dir, filename := clearRootPath(s.dir, clearRoot)
	f, err := os.Open(path.Join(dir, filename))
	if os.IsNotExist(err) {
		return mid.CodeErr{C: http.StatusNotFound}
	}
	if err != nil {
		return errors.Wrapf(err, "opening %s", filename)
	}
	defer f.Close()

	contentType, err := os.ReadFile(path.Join(dir, "content-type"))
	if err != nil {
		return errors.Wrap(err, "getting content type")
	}

	var tokenType common.Address
	if tokenTypeHex != "" {
		tokenType = common.HexToAddress(tokenTypeHex)
	}

	amount := new(big.Int)
	amount.SetString(amountStr, 10)
	if amount.Cmp(big0) < 1 {
		return fmt.Errorf("got amount = %s, should be 1 or greater", amount)
	}

	collateral := new(big.Int)
	collateral.SetString(collateralStr, 10)
	if collateral.Cmp(big0) < 1 {
		return fmt.Errorf("got collateral = %s, should be 1 or greater", collateral)
	}

	if err := s.checkPrice(tokenType, amount, collateral, clearRoot); err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "proposed payment rejected")}
	}

	revealDeadlineSecs, err := strconv.ParseInt(revealDeadlineSecsStr, 10, 64)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "parsing reveal deadline")}

	}
	revealDeadline := time.Unix(revealDeadlineSecs, 0)

	if time.Until(revealDeadline) < minRevealDur {
		return mid.CodeErr{C: http.StatusBadRequest, Err: fmt.Errorf("reveal deadline too soon: %s, require %s", time.Until(revealDeadline), minRevealDur)}
	}

	refundDeadlineSecs, err := strconv.ParseInt(refundDeadlineSecsStr, 10, 64)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "parsing refund deadline")}

	}
	refundDeadline := time.Unix(refundDeadlineSecs, 0)

	if refundDeadline.Sub(revealDeadline) > maxRefundDur {
		return mid.CodeErr{C: http.StatusBadRequest, Err: fmt.Errorf("refund deadline too much later after reveal deadline: %s, require %s", refundDeadline.Sub(revealDeadline), maxRefundDur)}
	}

	var key, transferID [32]byte

	if _, err := rand.Read(transferID[:]); err != nil {
		return errors.Wrap(err, "choosing transfer ID")
	}
	if _, err := rand.Read(key[:]); err != nil {
		return errors.Wrap(err, "choosing cipher key")
	}

	log.Printf("new transfer %x, clearRoot %x, payment %s/%s, collateral %s, deadlines %s/%s, key %x", transferID[:], clearRoot, amount, tokenType, collateral, revealDeadline, refundDeadline, key[:])

	w.Header().Set("X-Tredd-Transfer-Id", hex.EncodeToString(transferID[:]))
	w.Header().Set("Content-Type", string(contentType))

	tmpfile, err := os.CreateTemp("", "treddserve")
	if err != nil {
		return errors.Wrap(err, "creating response tempfile")
	}
	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)
	defer tmpfile.Close()

	cipherRoot, err := tredd.Serve(tmpfile, f, key)
	if err != nil {
		return errors.Wrap(err, "serving data")
	}
	var cipherRootBuf [32]byte
	copy(cipherRootBuf[:], cipherRoot)

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "closing response tempfile")
	}

	rec := &serverRecord{
		transferID:     transferID,
		tokenType:      tokenType,
		amount:         amount,
		collateral:     collateral,
		revealDeadline: revealDeadline,
		refundDeadline: refundDeadline,
		buyer:          buyer,
		key:            key,
		clearRoot:      clearRoot,
		cipherRoot:     cipherRootBuf,
	}

	if err := s.storeRecord(req.Context(), rec); err != nil {
		return errors.Wrap(err, "storing transfer record")
	}

	tmpfile, err = os.Open(tmpfilename)
	if err != nil {
		return errors.Wrap(err, "reopening response tempfile")
	}
	defer tmpfile.Close()

	_, err = io.Copy(w, tmpfile)
	return errors.Wrap(err, "writing response")
}

func (s *server) revealKey(w http.ResponseWriter, req *http.Request) error {
	var (
		transferIDHex   = req.FormValue("transferid")
		contractAddrHex = req.FormValue("contractaddr")
	)

	transferID, err := hex.DecodeString(transferIDHex)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "decoding transfer ID")}
	}

	contractAddr := common.HexToAddress(contractAddrHex)

	ctx := req.Context()
	rec, err := s.getRecord(ctx, transferID)
	if err != nil {
		return errors.Wrap(err, "finding transfer record")
	}
	rec.contractAddr = &contractAddr
	if err := s.storeRecord(ctx, rec); err != nil {
		return errors.Wrap(err, "updating transfer record")
	}

	con, receipt, err := tredd.RevealKey(
		ctx,
		s.client,
		time.Now(),
		s.seller,
		contractAddr,
		rec.key,
		rec.tokenType,
		rec.amount, rec.collateral,
		rec.revealDeadline, rec.refundDeadline,
		rec.clearRoot, rec.cipherRoot,
	)
	if err != nil {
		return mid.CodeErr{C: http.StatusBadRequest, Err: errors.Wrap(err, "constructing reveal-key transaction")}
	}

	log.Printf("revealed key in transaction %x", receipt.TxHash[:])

	s.queueClaimPaymentHelper(ctx, rec, con)

	return nil
}

func (s *server) getRecord(ctx context.Context, transferID []byte) (*serverRecord, error) {
	var rec serverRecord
	copy(rec.transferID[:], transferID)

	const q = `
    SELECT contract_addr, token_type, amount, collateral, reveal_deadline_secs, refund_deadline_secs, buyer, key, clear_root, cipher_root
    	FROM transfers
    	WHERE transfer_id = $1
	`

	var (
		contractAddr                           []byte
		revealDeadlineSecs, refundDeadlineSecs int64
		amount, collateral                     string
	)
	if err := s.db.QueryRowContext(ctx, q, transferID).Scan(&contractAddr, &rec.tokenType, &amount, &collateral, &revealDeadlineSecs, &refundDeadlineSecs, &rec.buyer, &rec.key, &rec.clearRoot, &rec.cipherRoot); err != nil {
		return nil, errors.Wrapf(err, "querying transfer record %x from db", transferID)
	}

	if len(contractAddr) > 0 {
		rec.contractAddr = new(common.Address)
		copy(rec.contractAddr[:], contractAddr)
	}

	rec.amount = new(big.Int)
	rec.amount.SetString(amount, 10)
	rec.collateral = new(big.Int)
	rec.collateral.SetString(collateral, 10)

	rec.revealDeadline = time.Unix(revealDeadlineSecs, 0)
	rec.refundDeadline = time.Unix(refundDeadlineSecs, 0)

	return &rec, nil
}

func (s *server) storeRecord(ctx context.Context, rec *serverRecord) error {
	const q = `
		INSERT OR REPLACE INTO transfers
			(transfer_id, contract_addr, token_type, amount, collateral, reveal_deadline_secs, refund_deadline_secs, buyer, key, clear_root, cipher_root)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	var contractAddr any
	if rec.contractAddr != nil {
		contractAddr = *rec.contractAddr
	}
	_, err := s.db.ExecContext(ctx, q, rec.transferID[:], contractAddr, rec.tokenType, rec.amount.String(), rec.collateral.String(), rec.revealDeadline.Unix(), rec.refundDeadline.Unix(), rec.buyer, rec.key, rec.clearRoot, rec.cipherRoot)
	return err
}

func (s *server) queueClaimPayment(ctx context.Context, transferID []byte) error {
	rec, err := s.getRecord(ctx, transferID)
	if err != nil {
		return errors.Wrap(err, "reading transfer record")
	}
	con := serverTreddABI.Instance(s.client, *rec.contractAddr)
	s.queueClaimPaymentHelper(ctx, rec, con)
	return nil
}

func (s *server) queueClaimPaymentHelper(ctx context.Context, rec *serverRecord, con *bind.BoundContract) {
	time.AfterFunc(time.Until(rec.refundDeadline), func() {
		tx, err := bind.Transact(con, s.seller, serverTreddABI.PackClaimPayment())
		if err != nil {
			log.Printf("ERROR claiming payment: %s", err)
			return
		}
		if _, err := bind.WaitMined(ctx, s.client, tx.Hash()); err != nil {
			log.Printf("ERROR awaiting claim-payment transaction: %s", err)
			return
		}
		if _, err := s.db.ExecContext(ctx, `DELETE FROM transfers WHERE transfer_id = $1`, rec.transferID); err != nil {
			log.Printf("ERROR deleting row from transfers table: %s", err)
		}
	})

	// TODO: set a timer that does tredd.ClaimPayment after the refund deadline
	// It should also delete the row from transfer_records.
}

func (s *server) checkPrice(tokenType common.Address, amount, collateral *big.Int, clearRoot [32]byte) error {
	// TODO: express seller preferences here (accepted currencies, per-item pricing, max collateral).
	return nil
}
