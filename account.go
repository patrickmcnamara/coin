package coin

import (
	"bytes"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
)

// Account is a coin account. It can make transactions on a ledger. It has a
// public key, which is used as the address of the account, and private key,
// which is used to sign transactions. Anyone with access to the private key
// will have access to the account.
type Account struct {
	PublicKey  PublicKey
	PrivateKey PrivateKey
}

// NewAccount returns a new account with a new public and private key.
func NewAccount() Account {
	pubKey, prvKey, _ := ed25519.GenerateKey(nil)
	return Account{
		PublicKey:  publicKeyConv(pubKey),
		PrivateKey: privateKeyConv(prvKey),
	}
}

// NewAccountFromSeed returns a new account with a public and private key
// generated from a given seed.
func NewAccountFromSeed(seed []byte) Account {
	hash := blake2b.Sum512(seed)
	pubKey, prvKey, _ := ed25519.GenerateKey(bytes.NewBuffer(hash[:]))
	return Account{
		PublicKey:  publicKeyConv(pubKey),
		PrivateKey: privateKeyConv(prvKey),
	}
}

// Sign signs data with the private key of the account.
func (acc Account) Sign(data []byte) Signature {
	return signatureConv(ed25519.Sign(acc.PrivateKey[:], data))
}

// Verify verifies signed data with the public key of the account.
func (acc Account) Verify(data []byte, sig Signature) bool {
	return ed25519.Verify(acc.PublicKey[:], data, sig[:])
}

// NewGenesisTransaction creates a new transaction where the account grants
// itself an amount of coin. This must be the first transaction in a ledger and
// will be invalid otherwise.
func (acc Account) NewGenesisTransaction(amount uint32) Transaction {
	return acc.NewTransaction(acc.PublicKey, amount, Signature{})
}

// NewTransaction creates a new transaction where an account send an amount of
// coin to another account, addressed by their respective public keys. It is
// signed by the private key of the sender. The current signature of the ledger
// that the transaction is to be added to is also required.
func (acc Account) NewTransaction(pubKey PublicKey, amount uint32, ledSig Signature) Transaction {
	trn := Transaction{
		From:   acc.PublicKey,
		To:     pubKey,
		Amount: amount,
	}
	trn.Sign(acc.PrivateKey, ledSig)
	return trn
}
