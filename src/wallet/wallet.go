package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"strings"

	"github.com/bmalcherek/goCoin/src/transaction"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
}

func InitWallet(password string) *Wallet {
	var sb strings.Builder
	sb.WriteString(password)
	for sb.Len() < 40 {
		sb.WriteString(password)
	}

	r := strings.NewReader(sb.String())
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		panic(err)
	}

	return &Wallet{
		PrivateKey: key,
	}
}

func (w *Wallet) MakeTransaction(address *ecdsa.PublicKey, amount float64, uTxOuts []*transaction.UnspentTxOut) *transaction.Transaction {
	uTxOutsForTx, err := transaction.FindUnspentTxOutsForTransaction(uTxOuts, &w.PrivateKey.PublicKey, amount)
	if err != nil {
		panic(err)
	}

	txIns := make([]*transaction.TxIn, len(uTxOutsForTx))
	for idx, tx := range uTxOutsForTx {
		txIn := tx.ConvertToUnsignedTxIn()
		txIn.Sign(w.PrivateKey)
		txIns[idx] = txIn
	}

	txOuts := transaction.GenerateTxOutsForTransaction(txIns, amount, &w.PrivateKey.PublicKey, address)

	t := &transaction.Transaction{
		TxIns:  txIns,
		TxOuts: txOuts,
	}
	t.Sign()

	return t
}

func (w *Wallet) GetBalance(uTxOuts []*transaction.UnspentTxOut) float64 {
	txs := transaction.FilterUnspentTxOuts(uTxOuts, &w.PrivateKey.PublicKey)
	balance := 0.0

	for _, tx := range txs {
		balance += tx.Amount
	}

	return balance
}

func (w *Wallet) String() string {
	var sb strings.Builder
	sb.WriteString(w.PrivateKey.D.String())

	return sb.String()
}
