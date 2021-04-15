package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"strings"
	"time"

	"github.com/bmalcherek/goCoin/src/constants"
	"github.com/bmalcherek/goCoin/src/crypto"
)

type Transaction struct {
	Id        string
	TxIns     []*TxIn
	TxOuts    []*TxOut
	Timestamp int64
}

func (t *Transaction) Sign() {
	var sb strings.Builder

	for _, tx := range t.TxOuts {
		sb.WriteString(tx.String())
	}

	for _, tx := range t.TxIns {
		sb.WriteString(tx.String())
	}
	sb.WriteString(fmt.Sprintf("%d", t.Timestamp))

	t.Id = crypto.Sha256(sb.String())
}

func CreateCoinbaseTransaction(address *ecdsa.PublicKey) *Transaction {
	t := &Transaction{
		TxOuts: []*TxOut{
			{
				Address: address,
				Amount:  constants.CoinbaseTransactionAmount,
			},
		},
		Timestamp: time.Now().UnixNano(),
	}

	t.Sign()

	return t
}

func HandleNewTransaction(t *Transaction, uTxOuts []*UnspentTxOut) []*UnspentTxOut {
	newUnspentTxOuts := lockUnspentTxOuts(t, uTxOuts)
	newUnspentTxOuts = mapUnspentTxOuts(t, newUnspentTxOuts)

	return newUnspentTxOuts
}

func lockUnspentTxOuts(t *Transaction, uTxOuts []*UnspentTxOut) []*UnspentTxOut {
	newUnspentTxOuts := []*UnspentTxOut{}

	// TODO improve finding this intersection
	for _, utx := range uTxOuts {
		locked := false
		for _, txIn := range t.TxIns {
			if utx.TransactionId == txIn.TransactionId && utx.TxOutIdx == txIn.TxOutIdx {
				locked = true
			}
		}

		if !locked {
			newUnspentTxOuts = append(newUnspentTxOuts, utx)
		}
	}

	return newUnspentTxOuts
}

func mapUnspentTxOuts(t *Transaction, uTxOuts []*UnspentTxOut) []*UnspentTxOut {
	newUnspentTxOuts := make([]*UnspentTxOut, len(t.TxOuts))

	for idx, tx := range t.TxOuts {
		newUnspentTxOuts[idx] = &UnspentTxOut{
			TransactionId: t.Id,
			TxOutIdx:      idx,
			Address:       tx.Address,
			Amount:        tx.Amount,
		}
	}

	return append(uTxOuts, newUnspentTxOuts...)
}

func (t *Transaction) String() string {
	var sb strings.Builder

	sb.WriteString(t.Id)

	for _, tx := range t.TxIns {
		sb.WriteString(tx.String())
	}

	for _, tx := range t.TxOuts {
		sb.WriteString(tx.String())
	}

	sb.WriteString(fmt.Sprintf("%d", t.Timestamp))

	return sb.String()
}
