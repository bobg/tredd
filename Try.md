# Trying Tredd

First, install the tredd binary.

```sh
$ go install github.com/bobg/tredd/cmd/tredd@latest
```

## Setting up an Ethereum test node

Tredd's Ethereum port requires a connection to an Ethereum node.
For local experimentation, use `geth` (go-ethereum) in developer mode,
which provides a single-node chain with a pre-funded account and instant block mining.

Install geth by following the instructions at https://geth.ethereum.org/docs/getting-started/installing-geth,
or via your system package manager.

Create a directory to hold tredd and geth files:

```sh
$ mkdir /path/to/dir
$ cd /path/to/dir
```

Start geth in developer mode with the HTTP RPC server enabled:

```sh
$ geth --dev --http --http.api eth,personal,web3 --datadir geth-dev-data 2>geth.log &
```

This starts geth listening on `http://127.0.0.1:8545` by default.
The `--dev` flag creates a pre-funded developer account,
auto-mines blocks, and persists chain data under `geth-dev-data/`.

## Creating seller and buyer accounts

Tredd identifies sellers and buyers by their Ethereum accounts.
Each account is represented by an encrypted keystore file.

Create a keystore directory and generate a keystore for the seller:

```sh
$ geth --datadir geth-dev-data account new
```

Enter a passphrase when prompted.
Geth will print the new account address and store the keystore file under
`geth-dev-data/keystore/`.
Note the address printed — this is the **seller address**.
Then run the same command again to create a **buyer account**,
and note its address separately.

For convenience, set shell variables for the addresses and keystore files:

```sh
$ SELLER_ADDR=0xYOUR_SELLER_ADDRESS
$ BUYER_ADDR=0xYOUR_BUYER_ADDRESS
$ SELLER_KEYFILE=geth-dev-data/keystore/YOUR_SELLER_KEYSTORE_FILENAME
$ BUYER_KEYFILE=geth-dev-data/keystore/YOUR_BUYER_KEYSTORE_FILENAME
```

## Funding the accounts

The `geth --dev` chain has a pre-funded developer account (the coinbase).
Attach a JavaScript console to the running node to send ETH to the seller and buyer:

```sh
$ geth attach geth-dev-data/geth.ipc
```

In the console:

```javascript
// Unlock the coinbase dev account (no passphrase needed in --dev mode)
eth.sendTransaction({from: eth.coinbase, to: "SELLER_ADDR", value: web3.toWei(10, "ether")})
eth.sendTransaction({from: eth.coinbase, to: "BUYER_ADDR",  value: web3.toWei(10, "ether")})
exit
```

Replace `SELLER_ADDR` and `BUYER_ADDR` with the hex addresses you noted above.
Both accounts now have 10 ETH to cover payments, collateral, and gas fees.

## Adding content to the server

Create a subdirectory for server content:

```sh
$ mkdir server-content
```

Add a file to the content tree:

```sh
$ tredd add -dir server-content /path/to/file
```

This command prints the content hash, for example:

```
added /path/to/file (content type text/plain; charset=utf-8) as a3b4c5...
```

Keep this hash — you will need it when issuing the buy request.

## Launching the Tredd server

Start the Tredd server, pointing it at the geth node and using the seller's keystore:

```sh
$ tredd serve \
    -dir server-content \
    -db server.db \
    -ethurl http://127.0.0.1:8545 \
    -keyfile $SELLER_KEYFILE \
    -passphrase YOUR_SELLER_PASSPHRASE
```

This creates `server.db` and logs the Tredd server's listen address
default `localhost:20544`.

## Buying content

In a separate shell, create a directory to hold downloaded content:

```sh
$ cd /path/to/dir
$ mkdir client-content
```

Issue a `get` request to the Tredd server:

```sh
$ tredd get \
    -hash CONTENT_HASH \
    -amount 1000000000000000000 \
    -collateral 1000000000000000000 \
    -reveal 15m \
    -refund 30m \
    -server http://localhost:20544 \
    -ethurl http://127.0.0.1:8545 \
    -seller $SELLER_ADDR \
    -keyfile $BUYER_KEYFILE \
    -passphrase YOUR_BUYER_PASSPHRASE \
    -dir client-content
```

Here:

- `CONTENT_HASH` is the hex hash reported by `tredd add`.
- `-amount` and `-collateral` are denominated in **wei**
  (1 ETH = 10¹⁸ wei; the example above proposes 1 ETH payment and 1 ETH collateral).
  Adjust to taste — even `1` wei works on a local test chain.
- `-reveal` is how long the seller has to publish the decryption key on-chain.
- `-refund` is how long after the reveal deadline the buyer has to claim a refund
  if the decrypted content is wrong.
  It must be at most 1 hour after the reveal deadline.
- `-seller` is the seller's Ethereum address (hex).
- To pay with an ERC20 token instead of ETH, add `-token TOKEN_CONTRACT_ADDRESS`.

## What happens under the hood

1. The client sends an HTTP request to the Tredd server proposing the payment terms.
2. The server encrypts the content with a fresh random key and streams back the
   ciphertext chunks interleaved with their clear-hashes.
   The client stores these and verifies they form a Merkle tree whose root
   matches `CONTENT_HASH`.
3. The client deploys a Tredd smart contract on Ethereum, funded with the proposed
   payment amount.
   The contract encodes the agreed terms (seller address, payment amount, collateral,
   Merkle roots, and deadlines).
4. The client notifies the server of the contract address.
5. The server calls `reveal()` on the contract, supplying the decryption key
   and the required collateral.
   This emits an on-chain event containing the key.
6. The client watches the chain, receives the key event, and decrypts the content.
   It verifies each decrypted chunk against its clear-hash.
   - If all chunks are correct, the transfer is complete.
     The decrypted file is written to `client-content/CONTENT_HASH`.
   - If a chunk is wrong, the client calls `refund()` on the contract with a
     Merkle proof of the bad chunk, reclaiming the payment plus the seller's collateral.
7. After the refund deadline passes, the server calls `claimPayment()` to collect
   the buyer's payment and reclaim its collateral.

## Utility subcommands

```sh
# Print the Tredd contract ABI
$ tredd abi

# Manually decrypt a stream of ciphertext chunks read from stdin
# (supply the 32-byte key as hex)
$ tredd decrypt -key KEYHEX < encrypted-file > decrypted-file
```