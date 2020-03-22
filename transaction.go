package coin

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
)

// Transaction is a coin transaction. An amount of coin is sent from one account
// to another. The sending account signs the transaction with its private key.
type Transaction struct {
	From      PublicKey `json:"from"`
	To        PublicKey `json:"to"`
	Amount    uint32    `json:"amount"`
	Signature Signature `json:"signature"`
}

// Sign returns the signature of a transaction given the private key of the
// sender and the signature of the ledger or bank.
func (trn Transaction) Sign(prvKey PrivateKey, currSig Signature) (sig Signature) {
	copy(sig[:], ed25519.Sign(prvKey[:], trn.Contract(currSig)))
	return
}

// Verify verifies the signature of the transaction with the public key of the
// sender and the signature of the ledger or bank.
func (trn Transaction) Verify(currSig Signature) bool {
	return ed25519.Verify(trn.From[:], trn.Contract(currSig), trn.Signature[:])
}

// Contract returns the bytes that account signs to create a transaction with
// the private key.
func (trn Transaction) Contract(currSig Signature) []byte {
	amount := make([]byte, 4)
	binary.LittleEndian.PutUint32(amount, trn.Amount)
	return bytes.Join([][]byte{trn.From[:], trn.To[:], amount, currSig[:]}, []byte{})
}
