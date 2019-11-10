package coin

import "bytes"

// Ledger is a list of transactions.
type Ledger struct {
	trns []Transaction
}

// NewLedger returns a new, empty ledger.
//
// You could also simply declare a Ledger instead.
func NewLedger() Ledger {
	return Ledger{}
}

// AddGenesisTransaction adds a genesis transaction to the ledger. Only a
// genesis transaction is valid as the first transaction of the ledger.
func (led *Ledger) AddGenesisTransaction(trn Transaction) error {
	if len(led.trns) > 0 {
		return ErrLedAlreadyGenesis
	}
	led.trns = append(led.trns, trn)
	return nil
}

// AddTransaction verifies that a transaction is valid and adds it to the ledger
// if it is.
func (led *Ledger) AddTransaction(trn Transaction) error {
	if len(led.trns) == 0 {
		return ErrLedNoGenesis
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
	trn, _ := led.LatestTransaction()
	return trn.Signature
}

// Transactions returns all transactions in the ledger.
func (led Ledger) Transactions() []Transaction {
	trns := make([]Transaction, len(led.trns))
	copy(trns, led.trns)
	return trns
}

// Do calls the given function on each transaction in the ledger, in order. If
// any of the calls return an error, Do will return that error immediately.
func (led Ledger) Do(f func(trn Transaction) error) error {
	for _, trn := range led.trns {
		if err := f(trn); err != nil {
			return err
		}
	}
	return nil
}

// TransactionsOf will return a slice of transactions in the ledger involving an
// account given its public key.
func (led Ledger) TransactionsOf(pubkey PublicKey) []Transaction {
	var trns []Transaction
	for _, trn := range led.trns {
		if bytes.Equal(trn.From[:], pubkey[:]) || bytes.Equal(trn.To[:], pubkey[:]) {
			trns = append(trns, trn)
		}
	}
	return trns
}

// GenesisTransaction returns the first transaction of the ledger.
func (led Ledger) GenesisTransaction() (trn Transaction, err error) {
	if len(led.trns) == 0 {
		err = ErrLedNoGenesis
		return
	}
	trn = led.trns[0]
	return
}

// LatestTransaction returns the latest transaction of the ledger.
func (led Ledger) LatestTransaction() (trn Transaction, err error) {
	if len(led.trns) == 0 {
		err = ErrLedNoGenesis
		return
	}
	trn = led.trns[len(led.trns)-1]
	return
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
