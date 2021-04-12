package transaction

import (
	"crypto/ecdsa"

	"github.com/bmalcherek/goCoin/src/constants"
)

type Transaction struct {
	TxOuts []*TxOut
}

func CreateCoinbaseTransaction(address *ecdsa.PublicKey) *Transaction {
	return &Transaction{
		TxOuts: []*TxOut{
			{
				Address: address,
				Amount:  constants.CoinbaseTransactionAmount,
			},
		},
	}
}
