package coin

import (
	"encoding/binary"

	"golang.org/x/crypto/ed25519"
)

// PublicKey is a public key and address of an account.
type PublicKey [ed25519.PublicKeySize]byte

// PrivateKey is a private key of an account.
type PrivateKey [ed25519.PrivateKeySize]byte

// Signature is a signature signed with the private key of an account.
type Signature [ed25519.SignatureSize]byte

// BurnAddress is a public key of an account that is inaccessible.
var BurnAddress PublicKey

func uint32ToBytes(num uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
}

func publicKeyConv(pk ed25519.PublicKey) (pubKey PublicKey) {
	copy(pubKey[:], pk)
	return pubKey
}

func privateKeyConv(pk ed25519.PrivateKey) (prvKey PrivateKey) {
	copy(prvKey[:], pk)
	return prvKey
}

func signatureConv(s []byte) (sig Signature) {
	copy(sig[:], s)
	return sig
}
