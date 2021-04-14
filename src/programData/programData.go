package programData

import (
	"sync"

	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/bmalcherek/goCoin/src/wallet"
)

type Store struct {
	Transactions  []*transaction.Transaction
	UnspentTxOuts []*transaction.UnspentTxOut
	Wallets       []*wallet.Wallet
	Lock          *sync.RWMutex
}
