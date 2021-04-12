package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"strings"
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

// func (w *Wallet) MakeTransaction(address *ecdsa.PublicKey) *transaction.Transaction {

// }

func (w *Wallet) String() string {
	var sb strings.Builder
	sb.WriteString(w.PrivateKey.D.String())

	return sb.String()
}
