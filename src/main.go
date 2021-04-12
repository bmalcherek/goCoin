package main

import (
	"log"

	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/bmalcherek/goCoin/src/wallet"
)

func main() {
	t := make([]*transaction.Transaction, 0)

	w1 := wallet.InitWallet("test")
	w2 := wallet.InitWallet("tesa")

	t = append(t, transaction.CreateCoinbaseTransaction(&w1.PrivateKey.PublicKey))
	t = append(t, transaction.CreateCoinbaseTransaction(&w2.PrivateKey.PublicKey))

	log.Println(w1.PrivateKey.PublicKey)
	log.Println(w1.PrivateKey)
	log.Println(w2.PrivateKey.PublicKey)

}
