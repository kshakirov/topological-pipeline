package main

import (
	"fmt"
	"net"
	"sync"
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
				Debug("TcpExternalChoice: received payload of bytes: ", "n", n, "content", buf[0])
				choice.BoxesChans[choice.currentIndex] <- buf[0] //for POC just
				choice.currentIndex = (choice.currentIndex + 1) % len(choice.BoxesChans)

			} else if err != nil {
				Debug("Can't read from connection")
				return
			}
		}
	}

	if len(choice.BoxesChans) == 0 {
		panic("TcpExternalChoice: requires at least one worker channel")
	}
	for _, workerInput := range choice.BoxesChans {
		if workerInput == nil {
			panic("TcpExternalChoice: worker channel must not be nil")
		}
	}
	go func() {
		Debug("TcpExternalChoice: Starting server on the host: ", "hostname", choice.Config.hostName)
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", choice.Config.hostName, choice.Config.port))
		if err != nil {
			Debug("TcpExternalChoice: Can't listen on the port ", "port", choice.Config.port) // handle error
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				Debug("TcpExternalChoice: Someting wrong with the connection ", "error", err)
			}
			Debug("TcpExternalChoice: Got Connection from a client")
			go handleConnection(conn)
		}
	}()

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

			Debug("TcpSplitBuffer: Can't open the connection to ", "host", smoother.Config.hostName, "port", smoother.Config.port)
		}
		for msg := range smoother.InChan {
			Debug("TcpSplitBuffer: received payload", "msg", msg)
			//			fmt.Fprint(conn, msg)
			conn.Write([]byte{msg})
			Debug("TcpSplitBuffer: emitted payload", "msg", msg)

		}
	}()
}

type TcpSplitNode struct {
	Nodes  []Node
	Buffer TcpSplitBuffer
}

func (node *TcpSplitNode) Process() {
	go func() {
		var workers sync.WaitGroup
		for _, worker := range node.Nodes {
			workers.Add(1)
			go func(worker Node) {
				defer workers.Done()
				for msg := range worker.InChan {
					result := worker.InBox.UserFuncB(msg)
					Debug("Box processed payload", "id", worker.ID, "result", result)
					node.Buffer.InChan <- result
				}
			}(worker)
		}

		workers.Wait()
		close(node.Buffer.InChan)
	}()
}
