package blockchain

import (
	"errors"

	"github.com/bmalcherek/goCoin/src/constants"
	"github.com/bmalcherek/goCoin/src/transaction"
)

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) GetLastBlock() (*Block, error) {
	bcLen := len(bc.Blocks)
	if bcLen == 0 {
		return nil, errors.New("blockchain is empty")
	}

	return bc.Blocks[bcLen-1], nil
}

func InitializeBlockchain() *Blockchain {

	genesisBlock := &Block{
		Index:        0,
		Hash:         "1a99af6e344146ab5d93ce075476f0fb0ba4f93e90c43434114b7a8a9adf5071",
		PreviousHash: "",
		Timestamp:    0,
		Transactions: []*transaction.Transaction{},
		Difficulty:   15,
	}

	bc := &Blockchain{
		Blocks: []*Block{genesisBlock},
	}

	GenerateBlock(bc, []*transaction.Transaction{})

	return bc
}

func (bc *Blockchain) GetDifficulty() int {
	last, err := bc.GetLastBlock()
	if err != nil {
		panic(err)
	}

	if last.Index%constants.DifficultyAdjustmentInterval == 0 && last.Index != 0 {
		return bc.getAdjustedDifficulty()
	}

	return last.Difficulty
}

func (bc *Blockchain) getAdjustedDifficulty() int {
	prevAdjBlock := bc.Blocks[len(bc.Blocks)-constants.DifficultyAdjustmentInterval-1]

	last, err := bc.GetLastBlock()
	if err != nil {
		panic(err)
	}

	timeExpected := int(constants.BlockGenerationInterval) * constants.DifficultyAdjustmentInterval
	timeTaken := last.Timestamp - prevAdjBlock.Timestamp
	if timeTaken < int64(timeExpected)/2 {
		return prevAdjBlock.Difficulty + 1
	} else if timeTaken > int64(timeExpected)*2 {
		return prevAdjBlock.Difficulty - 1
	} else {
		return prevAdjBlock.Difficulty
	}
}
