package miner

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bmalcherek/goCoin/src/blockchain.go"
	"github.com/bmalcherek/goCoin/src/crypto"
	"github.com/bmalcherek/goCoin/src/programData"
	"github.com/bmalcherek/goCoin/src/transaction"
)

func Mine(pd *programData.Store, bc *blockchain.Blockchain) {
	rand.Seed(time.Now().UnixNano())

	for {
		txs := []*transaction.Transaction{
			transaction.CreateCoinbaseTransaction(&pd.Wallets[0].PrivateKey.PublicKey),
		}
		txs = append(txs, pd.Transactions...)
		newBlock := blockchain.GenerateBlock(bc, txs)
		blockString := newBlock.String()
		difPrefix := getDifficultyPrefix(newBlock)
		nonce := rand.Int63()
		for {
			if HashMatchesDifficulty(nonce, blockString, difPrefix) {
				newBlock.Nonce = nonce
				break
			}
			nonce++
		}
	}
}

func getDifficultyPrefix(block *blockchain.Block) string {
	var sb strings.Builder
	for i := 0; i < block.Difficulty; i++ {
		sb.WriteString("0")
	}

	return sb.String()
}

func HashMatchesDifficulty(nonce int64, blockString, difPrefix string) bool {
	hash := crypto.Sha256Bin(fmt.Sprintf("%s%d", blockString, nonce))

	return strings.HasPrefix(hash, difPrefix)
}
