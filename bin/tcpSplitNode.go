package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type PrefixGeneratorFuncTCP func()
type SinkFuncTCP func()

type PrefixGeneratorTCP struct {
	Func PrefixGeneratorFuncTCP
}

func (generator *PrefixGeneratorTCP) Start() {
	go generator.Func()
}

type SinkTCP struct {
	Func   SinkFuncTCP
	InChan chan byte
}

func (sink *SinkTCP) Consume() {

	sink.Func()
}

type TcpExternalChoiceConfig struct {
	hostName string
	port     int
}

type TcpExternalChoice struct {
	currentIndex int
	InChan       chan byte
	BoxesChans   []chan byte
	Config       TcpExternalChoiceConfig
}

func (choice *TcpExternalChoice) WriteWithChoice() {

	handleConnection := func(conn net.Conn) {
		defer conn.Close()
		buf := make([]byte, 64)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				Debug("TcpExternalChoice received payload ", "n", n)
				choice.BoxesChans[choice.currentIndex] <- buf[0] //for POC just
				choice.currentIndex = (choice.currentIndex + 1) % len(choice.BoxesChans)

			} else if err != nil {
				Debug("Can't read from connection")
				return
			}
		}
	}

	if len(choice.BoxesChans) == 0 {
		panic("LocalExternalChoice requires at least one worker channel")
	}
	for _, workerInput := range choice.BoxesChans {
		if workerInput == nil {
			panic("LocalExternalChoice worker channel must not be nil")
		}
	}
	go func() {
		Debug("Starting server on the host: ", "hostname", choice.Config.hostName)
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", choice.Config.hostName, choice.Config.port))
		if err != nil {
			log.Fatalf("Can't listen on the port %d\n", choice.Config.port) // handle error
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
	}()

	// go func() {
	// 	defer func() {
	// 		for _, workerInput := range choice.BoxesChans {
	// 			close(workerInput)
	// 		}
	// 	}()

	// 	for msg := range choice.InChan {
	// 		Debug("LocalExternalChoice received payload", "msg", msg)
	// 		choice.BoxesChans[choice.currentIndex] <- msg
	// 		choice.currentIndex = (choice.currentIndex + 1) % len(choice.BoxesChans)
	// 	}
	// }()
}

type TcpSpliterBuffferConfig struct {
	hostName string
	port     int
}

type TcpSplitBuffer struct {
	InChan chan byte
	Config TcpSpliterBuffferConfig
	//	OutChan chan byte
}

func (smoother *TcpSplitBuffer) Interleave() {

	go func() {

		time.Sleep(time.Second * 5)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", smoother.Config.hostName, smoother.Config.port))

		if err != nil {

			log.Fatalf("Can't open the connection to %s  port %d\n", smoother.Config.hostName, smoother.Config.port)
		}
		for msg := range smoother.InChan {
			Debug("LocalSplitBuffer received payload", "msg", msg)
			fmt.Fprint(conn, msg)
			Debug("PrefixGenerator emitted payload", "msg", msg)

		}
	}()
}


