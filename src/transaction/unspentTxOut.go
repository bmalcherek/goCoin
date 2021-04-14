package transaction

import (
	"crypto/ecdsa"
	"errors"
)

type UnspentTxOut struct {
	TransactionId string
	TxOutIdx      int
	Address       *ecdsa.PublicKey
	Amount        float64
}

func FilterUnspentTxOuts(uTxOuts []*UnspentTxOut, address *ecdsa.PublicKey) []*UnspentTxOut {
	txs := []*UnspentTxOut{}

	for _, tx := range uTxOuts {
		if tx.Address.X.Cmp(address.X) == 0 && tx.Address.Y.Cmp(address.Y) == 0 {
			txs = append(txs, tx)
		}
	}

	return txs
}

func FindUnspentTxOutsForTransaction(uTxOuts []*UnspentTxOut, address *ecdsa.PublicKey, amount float64) ([]*UnspentTxOut, error) {
	txs := FilterUnspentTxOuts(uTxOuts, address)

	acc := 0.0
	selectedTxs := []*UnspentTxOut{}

	for _, tx := range txs {
		selectedTxs = append(selectedTxs, tx)
		acc += tx.Amount

		if acc >= amount {
			return selectedTxs, nil
		}
	}

	return nil, errors.New("not enough coins")
}

func (t *UnspentTxOut) ConvertToUnsignedTxIn() *TxIn {
	return &TxIn{
		TransactionId: t.TransactionId,
		TxOutIdx:      t.TxOutIdx,
		Amount:        t.Amount,
	}
}
