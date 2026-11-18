package main

import (
	"sync"
	"sync/atomic"
)

type PrefixGeneratorFuncTCP func()
type SinkFuncTCP func()
type BoxFuncTcp func(byte) byte

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

type TcpSplitNodeConfig struct {
	ChoiceConfig TcpExternalChoiceConfig
	BufferConfig TcpSpliterBuffferConfig
	Parallelism  int
}

type TcpExternalChoiceConfig struct {
	hostName string
	port     int
}

type TcpExternalChoice struct {
	currentIndex uint64
	InChan       chan Message
	BoxesChans   []chan Message
	Config       TcpExternalChoiceConfig
}

func (choice *TcpExternalChoice) WriteWithChoice() {

	if len(choice.BoxesChans) == 0 {
		panic("TcpExternalChoice: requires at least one worker channel")
	}
	for _, workerInput := range choice.BoxesChans {
		if workerInput == nil {
			panic("TcpExternalChoice: worker channel must not be nil")
		}
	}
	go func() {
		for _, c := range choice.BoxesChans {
			defer close(c)
		}
		for msg := range choice.InChan {
			nextId := atomic.AddUint64(&choice.currentIndex, 1) - 1
			choice.BoxesChans[nextId%uint64(len(choice.BoxesChans))] <- msg //for POC just

		}

	}()

}

type TcpSpliterBuffferConfig struct {
	hostName string
	port     int
}

type TcpSplitBuffer struct {
	InChan  chan OutBound
	OutChan chan OutBound //this one to node agent
	Config  TcpSpliterBuffferConfig
	//	OutChan chan byte
}

func (smoother *TcpSplitBuffer) Interleave() {

	go func() {
		defer close(smoother.OutChan)
		for msg := range smoother.InChan {
			smoother.OutChan <- msg

		}
	}()
}

type TcpSplitNode struct {
	Choice TcpExternalChoice
	Nodes  []Node
	Buffer TcpSplitBuffer
	Agent  NodeAgent
}

func (node *TcpSplitNode) Process() {
	go func() {
		var workers sync.WaitGroup
		for _, worker := range node.Nodes {
			workers.Add(1)
			go func(worker Node) {
				defer workers.Done()
				for msg := range worker.InChan {
					results := worker.InBox.UserFuncB(msg)
					for _, outbound := range results {
						node.Buffer.InChan <- outbound
					}
				}
			}(worker)
		}

		workers.Wait()
		close(node.Buffer.InChan)
	}()
}
