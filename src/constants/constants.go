package constants

import "time"

const (
	CoinbaseTransactionAmount    float64 = 50
	BlockGenerationInterval      int64   = int64(20 * time.Second)
	DifficultyAdjustmentInterval int     = 10
)
