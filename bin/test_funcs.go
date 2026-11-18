package main

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

func testComputeB(msg Message) []OutBound {
	Debug("testComputeB received message with RootId", "id", msg.RootId)
	completedID := msg.MessageId
	msg.ParentId = completedID
	msg.MessageId = uuid.New()
	msg.Payload[0] = msg.Payload[0] & 1

	return []OutBound{
		{
			Kind: COMPLETE,
			Complete: CompleteMessage{
				rootId:      msg.RootId,
				completedId: completedID,
				children:    []uuid.UUID{msg.MessageId},
			},
		},
		{Kind: DATA, Data: msg},
	}
}

func testPrefixGeneratorLoop(ctx context.Context) func(chan byte) {
	return func(out chan byte) {
		defer close(out)

		for {
			interval := rand.NormFloat64()*0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			timer := time.NewTimer(duration)

			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}

			value := byte(rand.IntN(256))
			out <- value
			Debug("PrefixGenerator emitted payload", "duration", duration, "value", value)
		}
	}
}

func testSinkFunc(input byte) {
	Debug("Sink received payload", "input", input)
}
