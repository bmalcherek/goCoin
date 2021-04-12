package main

import (
	"fmt"

	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/bmalcherek/goCoin/src/wallet"
)

func main() {
	t := make([]*transaction.Transaction, 0)
	uTxOuts := []*transaction.UnspentTxOut{}

	w1 := wallet.InitWallet("test")
	w2 := wallet.InitWallet("tesa")

	tx1 := transaction.CreateCoinbaseTransaction(&w1.PrivateKey.PublicKey)
	t = append(t, tx1)
	uTxOuts = transaction.HandleNewTransaction(tx1, uTxOuts)

	fmt.Printf("Wallet 1: %f\nWallet2: %f\n\n", w1.GetBalance(uTxOuts), w2.GetBalance(uTxOuts))

	tx2 := transaction.CreateCoinbaseTransaction(&w2.PrivateKey.PublicKey)
	t = append(t, tx2)
	uTxOuts = transaction.HandleNewTransaction(tx2, uTxOuts)

	fmt.Printf("Wallet 1: %f\nWallet2: %f\n\n", w1.GetBalance(uTxOuts), w2.GetBalance(uTxOuts))

	tx3 := w1.MakeTransaction(&w2.PrivateKey.PublicKey, 30.3123124541212, uTxOuts)
	t = append(t, tx3)
	uTxOuts = transaction.HandleNewTransaction(tx3, uTxOuts)

	fmt.Printf("Wallet 1: %f\nWallet2: %f\n\n", w1.GetBalance(uTxOuts), w2.GetBalance(uTxOuts))

	tx4 := w2.MakeTransaction(&w1.PrivateKey.PublicKey, 65.2131231274574, uTxOuts)
	t = append(t, tx4)
	uTxOuts = transaction.HandleNewTransaction(tx4, uTxOuts)

	fmt.Printf("Wallet 1: %f\nWallet2: %f\n\n", w1.GetBalance(uTxOuts), w2.GetBalance(uTxOuts))
}
