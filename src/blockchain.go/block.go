package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/bmalcherek/goCoin/src/crypto"
	"github.com/bmalcherek/goCoin/src/transaction"
)

type Block struct {
	Index        int
	Hash         string
	PreviousHash string
	Timestamp    int64
	Transactions []*transaction.Transaction
	Difficulty   int
	Nonce        int64
}

func GenerateBlock(bc *Blockchain, txs []*transaction.Transaction) *Block {
	lastBlock, err := bc.GetLastBlock()
	if err != nil {
		panic(err)
	}

	newBlock := &Block{
		Index:        lastBlock.Index + 1,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().UnixNano(),
		Transactions: txs,
		Difficulty:   bc.GetDifficulty(),
	}

	newBlock.calculateHash()

	return newBlock
}

func (b *Block) calculateHash() {
	var sb strings.Builder

	sb.WriteString(b.String())

	b.Hash = crypto.Sha256(sb.String())
}

func (b *Block) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d", b.Index))
	sb.WriteString(fmt.Sprintf("%d", b.Timestamp))
	sb.WriteString(fmt.Sprintf("%d", b.Difficulty))
	sb.WriteString(b.PreviousHash)

	for _, tx := range b.Transactions {
		sb.WriteString(tx.String())
	}

	return sb.String()
}
