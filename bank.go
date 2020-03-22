package coin

import (
	"bytes"
	"sync"
)

// Bank is a store of balances. It processes transactions but doesn't keep a
// transaction history.
type Bank struct {
	sig  Signature
	bals map[PublicKey]uint32
	lock sync.RWMutex
}

// NewBank creates a new bank given a genesis transaction.
func NewBank(trn Transaction) (bnk *Bank, err error) {
	// create bank
	bnk = new(Bank)
	bnk.bals = make(map[PublicKey]uint32)

	// amount must not be zero
	if trn.Amount == 0 {
		err = ErrTrnAmountZero
		return
	}
	// signature must be ok
	if !trn.Verify(bnk.Signature()) {
		err = ErrTrnBadSignature
		return
	}

	// lock/unlock
	bnk.lock.Lock()
	defer bnk.lock.Unlock()

	// do transaction
	bnk.bals[trn.To] = trn.Amount
	bnk.sig = trn.Signature
	return
}

// Signature returns the current signature of the bank. This is the signature
// of the latest transaction.
func (bnk *Bank) Signature() Signature {
	bnk.lock.RLock()
	defer bnk.lock.RUnlock()
	return bnk.sig
}

// Balance returns the balance of an account in this bank given its public key.
// If the public key has never been used, it returns 0.
func (bnk *Bank) Balance(pubKey PublicKey) uint32 {
	bnk.lock.RLock()
	defer bnk.lock.RUnlock()
	return bnk.bals[pubKey]
}

// Transaction validates and processes the given transaction in the bank. It
// updates the balances of the accounts and the bank's current signature.
func (bnk *Bank) Transaction(trn Transaction) error {
	// amount must not be zero
	if trn.Amount == 0 {
		return ErrTrnAmountZero
	}
	// sender and receiver must not be the same
	if bytes.Equal(trn.From[:], trn.To[:]) {
		return ErrTrnSameReceiver
	}
	// sender must have enough coin
	if trn.Amount > bnk.Balance(trn.From) {
		return ErrTrnAmountBalance
	}
	// signature must be ok
	if !trn.Verify(bnk.Signature()) {
		return ErrTrnBadSignature
	}

	// lock/unlock
	bnk.lock.Lock()
	defer bnk.lock.Unlock()

	// do transaction
	bnk.bals[trn.To] += trn.Amount
	bnk.bals[trn.From] -= trn.Amount
	bnk.sig = trn.Signature

	// cleanup potential zero balance
	if bnk.bals[trn.From] == 0 {
		delete(bnk.bals, trn.From)
	}

	return nil
}
