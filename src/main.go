package main

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"sync"
	"time"

	"github.com/bmalcherek/goCoin/src/programData"
	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/bmalcherek/goCoin/src/wallet"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

var (
	store *programData.Store
)

func main() {
	gob.Register(elliptic.P256())

	// store = &Store{
	// 	Transactions: []*transaction.Transaction{},
	// 	Lock:         &sync.RWMutex{},
	// }
	store = &programData.Store{
		Transactions:  []*transaction.Transaction{},
		UnspentTxOuts: []*transaction.UnspentTxOut{},
		Wallets: []*wallet.Wallet{
			wallet.InitWallet("test"),
			wallet.InitWallet("abc"),
		},
		Lock: &sync.RWMutex{},
	}

	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Println("Usage: Run './chat -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("Now run './chat -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging,
	// off by default. Otherwise, it uses rand.Reader.
	var r io.Reader
	if *debug {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

	if *dest == "" {
		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", handleStream)

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})
	} else {
		fmt.Println("This node's multiaddresses:")
		for _, la := range host.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
		if err != nil {
			panic(err)
		}

		// Create a buffered stream so that read and writes are non blocking.
		// rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		// go writeData(rw)
		// go readData(rw)
		txChan := make(chan *transaction.Transaction)
		go readTransaction(&s, txChan)
		go handleNewTransaction(txChan)
		go transactionPrinter()

		// Hang forever.
		select {}
	}
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	// rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// go readData(rw)
	// go writeData(rw)
	go sendTransaction(&s)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func sendTransaction(s *network.Stream) {
	i := 0
	for {
		// w1 := wallet.InitWallet("test")
		store.Lock.RLock()
		idx := mrand.Intn(len(store.Wallets))
		t := transaction.CreateCoinbaseTransaction(&store.Wallets[idx].PrivateKey.PublicKey)
		store.Lock.RUnlock()
		enc := gob.NewEncoder(*s)
		enc.Encode(t)
		i++
		fmt.Printf("Sent transaction %d\n", i)
		time.Sleep(1 * time.Second)
	}
}

func readTransaction(s *network.Stream, out chan<- *transaction.Transaction) {
	for {
		t := transaction.Transaction{}
		dec := gob.NewDecoder(*s)
		dec.Decode(&t)
		out <- &t
	}
}

func handleNewTransaction(in <-chan *transaction.Transaction) {
	for t := range in {
		store.Lock.Lock()
		store.Transactions = append(store.Transactions, t)
		store.UnspentTxOuts = transaction.HandleNewTransaction(t, store.UnspentTxOuts)
		fmt.Printf("Got new transaction %s\n", t.Id)
		store.Lock.Unlock()
	}
}

func transactionPrinter() {
	for {
		store.Lock.RLock()
		fmt.Printf("Wallet 1: %f\n", store.Wallets[0].GetBalance(store.UnspentTxOuts))
		fmt.Printf("Wallet 2: %f\n\n", store.Wallets[1].GetBalance(store.UnspentTxOuts))
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
