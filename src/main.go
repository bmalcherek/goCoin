package main

import (
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bmalcherek/goCoin/src/p2p"
	"github.com/bmalcherek/goCoin/src/programData"
	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/bmalcherek/goCoin/src/wallet"
)

var (
	store *programData.Store
)

func main() {
	gob.Register(elliptic.P256())

	store = &programData.Store{
		Transactions:  []*transaction.Transaction{},
		UnspentTxOuts: []*transaction.UnspentTxOut{},
		Wallets: []*wallet.Wallet{
			wallet.InitWallet("test"),
			wallet.InitWallet("abc"),
		},
		Lock: &sync.RWMutex{},
	}

	p2p.Setup()

	go handleNewTransaction()
	go sendTransaction()
	go transactionPrinter()

	select {}
}

func sendTransaction() {
	for {
		store.Lock.RLock()
		idx := rand.Intn(len(store.Wallets))
		t := transaction.CreateCoinbaseTransaction(&store.Wallets[idx].PrivateKey.PublicKey)
		store.Lock.RUnlock()

		tps := p2p.GetTransactionPubSub()
		tps.Publish(t)

		fmt.Printf("Peer \x1b[32m%s\x1b[0m sent transaction %s\n\n", tps.Self.String(), t.Id)
		sleep := rand.Int63n(10)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func handleNewTransaction() {
	for t := range p2p.GetTransactionPubSub().Transactions {
		store.Lock.Lock()
		store.Transactions = append(store.Transactions, t)
		store.UnspentTxOuts = transaction.HandleNewTransaction(t, store.UnspentTxOuts)
		fmt.Printf("\nGot new transaction %s\n", t.Id)
		store.Lock.Unlock()
	}
}

func transactionPrinter() {
	for {
		store.Lock.RLock()
		fmt.Printf("Wallet 1: %f\n", store.Wallets[0].GetBalance(store.UnspentTxOuts))
		fmt.Printf("Wallet 2: %f\n", store.Wallets[1].GetBalance(store.UnspentTxOuts))
		store.Lock.RUnlock()
		time.Sleep(5 * time.Second)
	}
}

func exampleTransactions() {
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
