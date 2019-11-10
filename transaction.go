package coin

import (
	"bytes"
	"crypto/ed25519"
)

// Transaction is a coin transaction. An amount of coin is sent from one account
// to another. The sending account signs the transaction with its private key.
type Transaction struct {
	From      PublicKey
	To        PublicKey
	Amount    uint32
	Signature Signature
}

// Sign signs the transaction with the private key of the sender and the current
// signature of the ledger.
func (trn *Transaction) Sign(prvKey PrivateKey, ledSig Signature) {
	trn.Signature = signatureConv(ed25519.Sign(prvKey[:], trn.Contract(ledSig)))
}

// Verify verifies the signature of the transaction with the public key of the
// sender and the current signature of the ledger.
func (trn Transaction) Verify(ledSig Signature) bool {
	return ed25519.Verify(trn.From[:], trn.Contract(ledSig), trn.Signature[:])
}

// Contract returns the bytes that account signs to create a transaction with
// the private key.
func (trn Transaction) Contract(ledSig Signature) []byte {
	return bytes.Join([][]byte{trn.From[:], trn.To[:], uint32ToBytes(trn.Amount), ledSig[:]}, []byte{})
}
