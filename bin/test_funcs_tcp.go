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
		time.Sleep(time.Second * 5)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d",hostname, port))

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


func handleConnection(conn net.Conn){
	defer conn.Close()
	buf:=make([]byte, 64,)
	for {
		n,err:= conn.Read(buf)
		if n > 0{
			Debug("read  bytes from the connection ", "n",n)
			
		
		}else if err !=nil{
			Debug("Can't read from connection")
			return 
		}
	}

	
}
func testSinkGeneratorLoopTCP(hostname string, port int) func(){

	return func(){
		Debug("Starting server on the host: ", "hostname", hostname)
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", hostname, port))
		if err != nil {
			log.Fatalf("Can't listen on the port %d\n", port)	// handle error
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				log.Fatalf("Someting wrong with the connection %v\n", err)
			}
			Debug("Got Connection from a client")
			go handleConnection(conn)
		}
	}
}
