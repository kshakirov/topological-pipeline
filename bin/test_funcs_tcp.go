package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

func testPrefixGeneratorLoopTCP(hostname string, port int) func() {
	return func() {

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d0",hostname, port))

		if err != nil {

			log.Fatalf("Can't open the connection to %s  port %d\n",  hostname, port)
		}
		for {
			value := byte(rand.IntN(256))
			interval := rand.NormFloat64()*0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			//	timer := time.NewTimer(duration)
			time.Sleep(duration)
			fmt.Fprint(conn, value)
			Debug("PrefixGenerator emitted payload", "duration", duration, "value", value)
		}
	}
}
