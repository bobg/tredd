# Trying Tredd

First, install the tredd binary.

```sh
$ go install github.com/bobg/tredd/cmd/tredd
```

Install the binaries from the github.com/chain/txvm package:

```sh
$ go install github.com/chain/txvm/cmd/...
```

Download and install txvmbcd, a minimal TxVM blockchain server:

```sh
$ go get -u github.com/bobg/txvmbcd
```

Create a directory to hold tredd and txvmbcd files and subdirectories.

```sh
$ mkdir /path/to/dir
```

Launch an instance of txvmbcd in that directory:

```sh
$ cd /path/to/dir
$ txvmbcd -db txvmbcd.db
```

This will create the file `txvmbcd.db` and produce log output in the shell,
including the listen address of the txvmbcd server,
and the hash of the initial block,
both of which you’ll need later.

The txvmbcd blockchain is empty.
It needs to be populated with some data before it can be used for Tredd.

In a separate shell,
cd to the directory and create a private/public keypair for an asset issuer:

```sh
$ cd /path/to/dir
$ ed25519 gen | tee issuer.prv | ed25519 pub >issuer.pub
```

Compute the ID of the default asset produced by this issuer:

```sh
$ assetid 1 `hex <issuer.pub` >asset-id
```

Create a private/public keypair for a Tredd seller of some information:

```sh
$ ed25519 gen | tee seller.prv | ed25519 pub >seller.pub
```

Also for a Tredd buyer:

```sh
$ ed25519 gen | tee buyer.prv | ed25519 pub >buyer.pub
```

Build and submit a transaction to the blockchain that issues 100 units of that asset,
sending 50 to the seller and 50 to the buyer:

```sh
$ tx build issue -blockchain BLOCKCHAINID -quorum 1 -prv `hex <issuer.prv` -pub `hex <issuer.pub` -amount 100 output -quorum 1 -pub `hex <seller.pub` -amount 50 -assetid `hex <asset-id` output -quorum 1 -pub `hex <buyer.pub` -amount 50 -assetid `hex <asset-id` | curl --data-binary @- http://LISTENADDR/submit
```

Here,
BLOCKCHAINID is the hash of the initial block and LISTENADDR is the txvmbcd listen address,
both reported when txvmbcd started.

Now the blockchain is populated. It’s time to populate the Tredd server with some content.

First, make a subdirectory for server content:

```sh
$ mkdir server-content
```

Now choose some file and add it to that content tree:

```sh
$ tredd add -dir server-content /path/to/file
```

This command will report the hash of the added content. You will need that in a moment.

With content in the server tree, it’s time to launch the server:

```sh
$ tredd serve -dir server-content -db server.db -prv seller.prv -url http://LISTENADDR
```

Here, LISTENADDR is the address of the txvmbcd server, still running in another shell.

This will create the file `seller.db` and produce log output in the shell,
including the listen address of the Tredd server, which you’ll need in the next step.

With the txvmbcd and Tredd servers running, it’s time to request and pay for some content.

In a third shell,
cd to the directory and create a subdirectory to hold content retrieved from the Tredd server:

```sh
$ cd /path/to/dir
$ mkdir client-content
```

Finally, issue a “get” request to the Tredd server:

```sh
$ tredd get -hash HASH -amount 1 -asset `hex <asset-id` -reveal 15m -refund 15m -db client.db -prv buyer.prv -server http://TREDDLISTEN -bcurl http://TXVMBCDLISTEN -dir client-content
```

Here,
HASH is the hash of the desired content,
reported above by `tredd add`,
TREDDLISTEN is the address of the Tredd server,
and TXVMBCDLISTEN is the address of the txvmbcd server.

This command will send a request to the Tredd server proposing payment of 1 unit of our defined asset in exchange for the content identified by HASH.
The proposal also includes a deadline of 15 minutes for the server to publish its “reveal-key” transaction on the blockchain,
and another 15 minutes for the client to claim a refund if warranted.

The server accepts the proposal,
chooses a unique transfer ID and a random encryption key,
and responds with an encrypted copy of the desired content as chunks interleaved with each chunk’s “clear hash”
(the hash it should have after decryption).
The client stores these chunks and hashes,
and double-checks that the stream of hashes produces a Merkle tree whose root is the HASH given above.

Now the client constructs a partial TxVM transaction and sends it to the server.
The transaction includes payment and a call to the Tredd contract that enforces the agreed terms of the transfer.

The server completes the transaction by adding its collateral payment and revealing the decryption key.
It publishes the completed transaction to the blockchain and schedules a task for after the “refund deadline” to claim its payment.

The client,
observing the blockchain,
sees the completed transaction and parses out the server’s key.
It uses that to decrypt the content it received earlier,
checking along the way that each decrypted chunk has the hash it’s supposed to.
If it does,
the transfer is complete.
If some chunk doesn’t have the right hash,
it constructs and publishes a “claim-refund” transaction to the blockchain.

After a successful transfer,
the retrieved content is in the file client-content/HASH.
