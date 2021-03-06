package coin

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
)

// Account is a coin account. It can make transactions on a ledger. It has a
// public key, which is used as the address of the account, and private key,
// which is used to sign transactions. Anyone with access to the private key
// will have access to the account.
type Account struct {
	PublicKey  PublicKey  `json:"publicKey"`
	PrivateKey PrivateKey `json:"privateKey"`
}

// NewAccount returns a new account with a new public and private key.
func NewAccount() (acc Account) {
	// generate random seed
	seed := make([]byte, 32)
	rand.Read(seed)

	// generate account from seed
	acc, _ = NewAccountFromSeed(seed)
	return
}

// NewAccountFromSeed returns a new account with a public and private key
// generated from a given seed.
func NewAccountFromSeed(seed []byte) (acc Account, err error) {
	// seed must have a length of at least 32
	if len(seed) < 32 {
		err = ErrAccShortSeed
		return
	}

	// generate account from seed
	pubKey, priKey, _ := ed25519.GenerateKey(bytes.NewBuffer(seed))
	copy(acc.PublicKey[:], pubKey)
	copy(acc.PrivateKey[:], priKey)

	return
}

// Sign signs data with the private key of the account.
func (acc Account) Sign(data []byte) (sig Signature) {
	copy(sig[:], ed25519.Sign(acc.PrivateKey[:], data))
	return
}

// Verify verifies signed data with the public key of the account.
func (acc Account) Verify(data []byte, sig Signature) bool {
	return ed25519.Verify(acc.PublicKey[:], data, sig[:])
}

// NewGenesisTransaction creates a new transaction where the account grants
// itself an amount of coin. This must be the first transaction in a ledger
// or bank and will be invalid otherwise.
func (acc Account) NewGenesisTransaction(amount uint32) (trn Transaction) {
	trn = acc.NewTransaction(acc.PublicKey, amount, Signature{})
	return
}

// NewTransaction creates a new transaction where an account send an amount of
// coin to another account, addressed by their respective public keys. It is
// signed by the private key of the sender. The signature of the ledger or bank
// that the transaction is to be added to is also required.
func (acc Account) NewTransaction(pubKey PublicKey, amount uint32, currSig Signature) (trn Transaction) {
	trn.To = pubKey
	trn.From = acc.PublicKey
	trn.Amount = amount
	trn.Signature = trn.Sign(acc.PrivateKey, currSig)
	return
}
