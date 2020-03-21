package coin

import "crypto/ed25519"

// PublicKey is a public key and address of an account.
type PublicKey [ed25519.PublicKeySize]byte

// PrivateKey is a private key of an account.
type PrivateKey [ed25519.PrivateKeySize]byte

// Signature is a signature signed with the private key of an account.
type Signature [ed25519.SignatureSize]byte

// BurnAddress is a public key of an account that is inaccessible.
var BurnAddress PublicKey
