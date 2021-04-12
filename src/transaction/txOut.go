package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"strings"
)

type TxOut struct {
	Address *ecdsa.PublicKey
	Amount  float64
}

func GenerateTxOutsForTransaction(txIns []*TxIn, amount float64, sender, recipient *ecdsa.PublicKey) []*TxOut {
	acc := 0.0
	for _, tx := range txIns {
		acc += tx.Amount
	}

	txOuts := make([]*TxOut, 2)
	txOuts[0] = &TxOut{
		Address: recipient,
		Amount:  amount,
	}

	txOuts[1] = &TxOut{
		Address: sender,
		Amount:  acc - amount,
	}

	return txOuts
}

func (t *TxOut) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%x", t.Address))
	sb.WriteString(fmt.Sprintf("%f", t.Amount))

	return sb.String()
}
