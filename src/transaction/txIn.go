package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/bmalcherek/goCoin/src/crypto"
)

type TxIn struct {
	TransactionId string
	TxOutIdx      int
	Signature     string
	Amount        float64
}

func (t *TxIn) Sign(key *ecdsa.PrivateKey) {
	hash := crypto.Sha256(t.String())
	sig, err := ecdsa.SignASN1(rand.Reader, key, []byte(hash))
	if err != nil {
		panic(err)
	}

	t.Signature = string(sig)
}

func (t *TxIn) String() string {
	var sb strings.Builder

	sb.WriteString(t.TransactionId)
	sb.WriteString(fmt.Sprintf("%d", t.TxOutIdx))
	sb.WriteString(fmt.Sprintf("%f", t.Amount))
	sb.WriteString(t.Signature)

	return sb.String()
}
