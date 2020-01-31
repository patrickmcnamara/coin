package coin

import (
	"bytes"
	"sync"
)

// Ledger is a list of transactions.
type Ledger struct {
	trns []Transaction
	bals map[PublicKey]uint32
	lock sync.RWMutex
}

// NewLedger returns a new, empty ledger.
func NewLedger() *Ledger {
	return &Ledger{bals: make(map[PublicKey]uint32)}
}

// AddGenesisTransaction adds a genesis transaction to the ledger. Only a
// genesis transaction is valid as the first transaction of the ledger.
func (led *Ledger) AddGenesisTransaction(trn Transaction) error {
	// ledger should not have a genesis transaction
	if len(led.trns) > 0 {
		return ErrLedAlreadyGenesis
	}

	// amount must not be zero
	if trn.Amount == 0 {
		return ErrTrnAmountZero
	}

	// signature must be ok
	if !trn.Verify(led.Signature()) {
		return ErrTrnBadSignature
	}

	// update internal state
	led.lock.Lock()
	defer led.lock.Unlock()
	led.trns = append(led.trns, trn)
	led.bals[trn.To] += trn.Amount

	return nil
}

// AddTransaction verifies that a transaction is valid and adds it to the ledger
// if it is.
func (led *Ledger) AddTransaction(trn Transaction) error {
	// ledger should have a genesis transaction
	if len(led.trns) == 0 {
		return ErrLedNoGenesis
	}

	// amount must not be zero
	if trn.Amount == 0 {
		return ErrTrnAmountZero
	}

	// sender and receiver must not be the same
	if bytes.Equal(trn.From[:], trn.To[:]) {
		return ErrTrnSameReceiver
	}

	// sender must have enough coin
	if trn.Amount > led.BalanceOf(trn.From) {
		return ErrTrnAmountBalance
	}

	// signature must be ok
	if !trn.Verify(led.Signature()) {
		return ErrTrnBadSignature
	}

	// update internal state
	led.lock.Lock()
	defer led.lock.Unlock()
	led.trns = append(led.trns, trn)
	led.bals[trn.To] += trn.Amount
	led.bals[trn.From] -= trn.Amount

	return nil
}

// Signature returns the current signature of the ledger. This is the signature
// of the latest transaction.
func (led *Ledger) Signature() Signature {
	led.lock.RLock()
	defer led.lock.RUnlock()
	trn, _ := led.LatestTransaction()
	return trn.Signature
}

// Transactions returns all transactions in the ledger.
func (led *Ledger) Transactions() []Transaction {
	trns := make([]Transaction, len(led.trns))
	copy(trns, led.trns)
	return trns
}

// Do calls the given function on each transaction in the ledger, in order. If
// any of the calls return an error, Do will return that error immediately.
func (led *Ledger) Do(f func(trn Transaction) error) error {
	led.lock.RLock()
	defer led.lock.RUnlock()
	for _, trn := range led.trns {
		if err := f(trn); err != nil {
			return err
		}
	}
	return nil
}

// TransactionsOf will return a slice of transactions in the ledger involving an
// account given its public key.
func (led *Ledger) TransactionsOf(pubkey PublicKey) []Transaction {
	led.lock.RLock()
	defer led.lock.RUnlock()
	var trns []Transaction
	for _, trn := range led.trns {
		if bytes.Equal(trn.From[:], pubkey[:]) || bytes.Equal(trn.To[:], pubkey[:]) {
			trns = append(trns, trn)
		}
	}
	return trns
}

// GenesisTransaction returns the first transaction of the ledger.
func (led *Ledger) GenesisTransaction() (trn Transaction, err error) {
	// genesis transaction must exist
	if len(led.trns) == 0 {
		err = ErrLedNoGenesis
		return
	}

	// genesis transaction
	led.lock.RLock()
	defer led.lock.RUnlock()
	trn = led.trns[0]

	return
}

// LatestTransaction returns the latest transaction of the ledger.
func (led *Ledger) LatestTransaction() (trn Transaction, err error) {
	// genesis transaction must exist
	if len(led.trns) == 0 {
		err = ErrLedNoGenesis
		return
	}

	// latest transaction
	led.lock.RLock()
	defer led.lock.RUnlock()
	trn = led.trns[len(led.trns)-1]

	return
}

// Balances returns a map of every public key in the ledger and the balances of
// the accounts associated with those public keys.
func (led *Ledger) Balances() map[PublicKey]uint32 {
	led.lock.RLock()
	defer led.lock.RUnlock()
	bals := make(map[PublicKey]uint32)
	for pubKey, amt := range led.bals {
		bals[pubKey] = amt
	}
	return bals
}

// BalanceOf returns the balance of an account in the ledger given its public
// key. If the account is not in the ledger, it will return 0.
func (led *Ledger) BalanceOf(pubKey PublicKey) uint32 {
	led.lock.RLock()
	defer led.lock.RUnlock()
	return led.bals[pubKey]
}

// Size returns the number of transactions in the ledger.
func (led *Ledger) Size() uint64 {
	return uint64(len(led.trns))
}
