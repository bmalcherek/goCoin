package p2p

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/bmalcherek/goCoin/src/transaction"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// BufSize is the number of incoming messages to buffer for each topic.
const BufSize = 128

type TransactionPubSub struct {
	Transactions chan *transaction.Transaction

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	Self peer.ID
}

func (tps *TransactionPubSub) readLoop() {
	for {
		msg, err := tps.sub.Next(tps.ctx)
		if err != nil {
			close(tps.Transactions)
			return
		}

		// if msg.ReceivedFrom == tps.self {
		// 	continue
		// }

		t := new(transaction.Transaction)
		dec := gob.NewDecoder(bytes.NewReader(msg.Data))
		if err = dec.Decode(&t); err != nil {
			panic(err)
		}

		tps.Transactions <- t
	}
}

func (tps *TransactionPubSub) Publish(t *transaction.Transaction) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(t); err != nil {
		return err
	}

	return tps.topic.Publish(tps.ctx, b.Bytes())
}

func (tps *TransactionPubSub) ListPeers() []peer.ID {
	return tps.ps.ListPeers("transactions")
}

func JoinTransactionPubSub(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID) (*TransactionPubSub, error) {
	topic, err := ps.Join("transactions")
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	tps := &TransactionPubSub{
		Transactions: make(chan *transaction.Transaction),

		ctx:   ctx,
		ps:    ps,
		topic: topic,
		sub:   sub,

		Self: selfID,
	}

	go tps.readLoop()
	return tps, nil
}
