package coin

import "errors"

// Account errors
//
// These indicate why an account could not be generated.
var (
	ErrAccShortSeed = errors.New("seed must have a length greater than 32")
)

// Transaction errors
//
// These indicate why a transaction was invalid.
var (
	ErrTrnAmountZero    = errors.New("amount cannot be zero")
	ErrTrnAmountBalance = errors.New("sender balance too low")
	ErrTrnSameReceiver  = errors.New("receiver cannot be same as sender")
	ErrTrnBadSignature  = errors.New("signature cannot be validated")
)

// Ledger errors
//
// These indicate why a transaction could not be added to the ledger.
var (
	ErrLedAlreadyGenesis = errors.New("genesis transaction already in ledger")
	ErrLedNoGenesis      = errors.New("no genesis transaction in ledger")
)
