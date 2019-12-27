# Coin

![latest version](https://img.shields.io/github/v/tag/patrickmcnamara/coin?label=latest%20version)
![last commit](https://img.shields.io/github/last-commit/patrickmcnamara/coin)
![top language](https://img.shields.io/github/languages/top/patrickmcnamara/coin)
![license](https://img.shields.io/github/license/patrickmcnamara/coin?label=license)

This is a basic cryptocurrency library, made for fun.
It has no networking capabilities.
It has no consensus mechanism.
It has no CLI.

## How it works

### Account

A public and private key.
The public key is used as the "address" of the account.
The private key is used to sign transactions to prove that they came from the account.
Anyone with the private key has access to the account.

### Transaction

An object containing the public key of the sender account, the public key of the receiver account, the amount of coin to be sent, and the signature from the sender account.
The data signed in a transaction includes the signature of the previous transaction.
This is used by the ledger to chain transactions together.

A special case is a genesis transaction. In this type of transaction, the public key of the sender account and receiver account are the same.
The sending account grants itself an amount of coin. This is the first transaction in a ledger.

### Ledger

A list of transactions.
A genesis transaction is first before all other transactions and there is only one.
Transactions are verified when they are added.

## Installation

Run `go get -u github.com/patrickmcnamara/coin`.

## License

This project is licensed under the GNU Lesser General Public License v3.0.
