package main

import (
	"context"
	"math/rand/v2"
	"time"
)

func testComputeB(input byte) byte {
	Debug("testComputeB received payload", "input", input)
	return input & 1
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
