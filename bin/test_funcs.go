package main

import (
	"log"

	"math/rand/v2"
	"time"
)

func heavyCompute(in Set) Set {
	val := in.(int)
	log.Printf("Func heavyCompute input=%d\n", val)
	// Искусственный разброс по времени, чтобы проверить FIFO на выходе
	if val%2 == 0 {
		time.Sleep(50 * time.Millisecond) // Четные — тугодумы
	} else {
		time.Sleep(10 * time.Millisecond) // Нечетные — шустрики
	}
	return val

}

func testCompute(i Set) Set {
	log.Printf("Func TestCompute input= %d\n", i)
	val := i.(int)
	return val * 20
}

func testComputeB(i byte) byte {
	log.Printf("Func TestCompute input= %d\n", i)

	return i & 1
}

func testPrefixGeneratorFunc(i byte) byte {
	log.Printf("Func test prefix generator input= %d\n", i)
	return i | 1
}

func testPrefixGeneratorLoop() func(chan byte) {
	return func(outChan chan byte) {

		defer close(outChan)
		for  {
			interval:=   rand.NormFloat64() * 0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			val := byte(rand.IntN(256))
			time.Sleep(duration)
			outChan <- val
		}

	}
}

func testSinkFunc(b byte) {
	log.Printf("Sink: received %d\n", b)
}
