package transaction

import "crypto/ecdsa"

type TxOut struct {
	Address *ecdsa.PublicKey
	Amount  float64
}
