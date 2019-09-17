package coin

import (
	"bytes"
	"errors"
)

// Ledger is a list of transactions.
type Ledger struct {
	trns []Transaction
}

// NewLedger returns a new, empty ledger.
func NewLedger() Ledger {
	return Ledger{make([]Transaction, 0)}
}

// AddGenesisTransaction adds a genesis transaction to the ledger. Only a
// genesis transaction is valid as the first transaction of the ledger.
func (led *Ledger) AddGenesisTransaction(trn Transaction) error {
	if len(led.trns) > 0 {
		return errors.New("add genesis transaction: genesis transaction already in ledger")
	}
	led.trns = append(led.trns, trn)
	return nil
}

// AddTransaction verifies that a transaction is valid and adds it to the ledger
// if it is.
func (led *Ledger) AddTransaction(trn Transaction) error {
	if len(led.trns) == 0 {
		return errors.New("add transaction: no genesis transaction in ledger")
	}
	if err := trn.Check(*led); err != nil {
		return err
	}
	led.trns = append(led.trns, trn)
	return nil
}

// Signature returns the current signature of the ledger. This is the signature
// of the latest transaction.
func (led Ledger) Signature() Signature {
	return led.LatestTransaction().Signature
}

// Verify verifies the signature of every transaction in the ledger, thus
// checking if the ledger is valid. This is only required when importing a
// ledger as transactions are verified as they are added to the ledger.
func (led Ledger) Verify() bool {
	var currSig Signature
	for _, trn := range led.trns {
		if !trn.Verify(currSig) {
			return false
		}
		currSig = trn.Signature
	}
	return true
}

// Transactions returns all transactions in the ledger.
func (led Ledger) Transactions() []Transaction {
	return []Transaction(led.trns)
}

// Do calls the given function on each transaction in the ledger, in order. If
// any of the calls return an error, Do will return that error immediately.
func (led Ledger) Do(f func(trn Transaction) error) error {
	for _, trn := range led.trns {
		err := f(trn)
		if err != nil {
			return err
		}
	}
	return nil
}

// TransactionsOf will return a slice of transactions in the ledger involving an
// account given its public key.
func (led Ledger) TransactionsOf(pubkey PublicKey) []Transaction {
	trns := make([]Transaction, 0)
	for _, trn := range led.trns {
		if bytes.Equal(trn.From[:], pubkey[:]) || bytes.Equal(trn.To[:], pubkey[:]) {
			trns = append(trns, trn)
		}
	}
	return trns
}

// GenesisTransaction returns the first transaction of the ledger.
func (led Ledger) GenesisTransaction() Transaction {
	if len(led.trns) == 0 {
		return Transaction{}
	}
	return led.trns[0]
}

// LatestTransaction returns the latest transaction of the ledger.
func (led Ledger) LatestTransaction() Transaction {
	if len(led.trns) == 0 {
		return Transaction{}
	}
	return led.trns[len(led.trns)-1]
}

// Balances returns a map of every public key in the ledger and the balances of
// the accounts associated with those public keys.
func (led Ledger) Balances() map[PublicKey]uint32 {
	bals := make(map[PublicKey]uint32)
	for i, trn := range led.trns {
		bals[trn.To] += trn.Amount
		if i != 0 {
			bals[trn.From] -= trn.Amount
		}
	}
	return bals
}

// BalanceOf returns the balance of an account in the ledger given its public
// key. If the account is not in the ledger, it will return 0.
func (led Ledger) BalanceOf(pubKey PublicKey) uint32 {
	var bal uint32
	for _, trn := range led.trns {
		if bytes.Equal(trn.To[:], pubKey[:]) {
			bal += trn.Amount
		} else if bytes.Equal(trn.From[:], pubKey[:]) {
			bal -= trn.Amount
		}
	}
	return bal
}

// Size returns the number of transactions in the ledger.
func (led Ledger) Size() uint64 {
	return uint64(len(led.trns))
}
