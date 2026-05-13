# TREDD - Trustless escrow for digital data

This is Tredd,
a software library that allows a buyer and a seller of some information to exchange payment for data securely.

It relies on
[the Ethereum blockchain](https://ethereum.org/)
and includes a demonstration client and server
(in
[cmd/tredd](https://github.com/bobg/tredd/tree/master/cmd/tredd)).

- A buyer sends a request for some content to a seller.
- The seller responds with an encrypted copy of the content.
- The buyer publishes a payment-proposal contract to the blockchain.
- The seller reveals the decryption key via the contract, together with a collateral payment.
- The buyer verifies the content and decryption key, releasing payment and returning the collateral to the seller.
- If the content or decryption key cannot be verified, the buyer recovers the proposed payment and keeps the collateral.

For more information,
see
[this detailed explanation of Tredd’s design and operation](https://medium.com/@bob.glickstein/tredd-trustless-escrow-for-digital-data-1ce5a0018e95).

(There is also [an appendix](https://medium.com/@bob.glickstein/tredd-appendix-frontier-sets-84e60bd19d59).)

For the motivation behind Tredd, see [Why Tredd](Why.md).

For step-by-step instructions for running the Tredd server and client, see [Trying Tredd](Try.md).

## Building Tredd

You must have [Go](https://go.dev/) installed.

```sh
go tool task build
```
