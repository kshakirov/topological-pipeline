package main

import (
	"fmt"
	"math/rand/v2"
	"net"
	"time"
)

func testPrefixGeneratorLoopTCP(hostname string, port int) func() {
	return func() {
		time.Sleep(time.Second * 5)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))

		if err != nil {

			Debug("PrefixGenerator: Can't open the connection to  ", "hostname", hostname, "port", port)
		}
		for {
			value := []byte{byte(rand.IntN(256))}
			interval := rand.NormFloat64()*0.5 + 2
			duration := time.Duration(interval * float64(time.Second))
			//	timer := time.NewTimer(duration)
			time.Sleep(duration)
			_, err := conn.Write(value)
			if err != nil {
				Debug("PrefixGenerator: cannot write byte", "value", value[0])
			}
			Debug("PrefixGenerator: emitted payload", "duration", duration, "value", value)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 64)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			Debug("Sink: read  bytes from the connection ", "n", n, "byte", buf[0])

		} else if err != nil {
			Debug(" Sink: Can't read from connection")
			return
		}
	}

}
func testSinkGeneratorLoopTCP(hostname string, port int) func() {

	return func() {
		Debug("Sink: Starting server on the host: ", "hostname", hostname)
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", hostname, port))
		if err != nil {
			Debug("Sink: Can't listen on the ", "port", port) // handle error
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				Debug("Sink: Someting wrong with the connection ", "error", err)
			}
			Debug("Got Connection from a client")
			go handleConnection(conn)
		}
	}
}
