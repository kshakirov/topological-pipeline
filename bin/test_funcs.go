package main

import (
	"context"

	"math/rand/v2"
	"time"
)

func heavyCompute(in Set) Set {
	val := in.(int)
	Debug("Func heavyCompute input = ","input", val)
	// Искусственный разброс по времени, чтобы проверить FIFO на выходе
	if val%2 == 0 {
		time.Sleep(50 * time.Millisecond) // Четные — тугодумы
	} else {
		time.Sleep(10 * time.Millisecond) // Нечетные — шустрики
	}
	return val

}

func testCompute(i Set) Set {
	Debug("Func TestCompute input= ", "input", i)
	val := i.(int)
	return val * 20
}

func testComputeB(i byte) byte {
	Debug("Func TestCompute input= ", "input", i)

	return i & 1
}

func testPrefixGeneratorFunc(i byte) byte {
	Debug("Func test prefix generator input= ", "input", i)
	return i | 1
}

func testPrefixGeneratorLoop(ctx context.Context) func(chan byte) {
	return func(outChan chan byte) {

		defer close(outChan)
		for {
			interval := rand.NormFloat64()*0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			val := byte(rand.IntN(256))
			Debug("Sleeping for ","duration", duration)
			time.Sleep(duration)

			select {
			case <-ctx.Done():

				return
			case outChan <- val:
			}

		}

	}
}

func testSinkFunc(b byte) {
	Debug("Sink: received ", "input", b)
}
