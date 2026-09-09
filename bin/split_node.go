package main

import "sync"

type BoxFuncB func(in byte) byte

type BoxB struct {
	ID        string
	UserFuncB BoxFuncB
}

type PrefixGeneratorFunc func(out chan byte)
type SinkFunc func(byte)

func NewBoxB(id string, fn BoxFuncB) *BoxB {
	return &BoxB{ID: id, UserFuncB: fn}
}

type LocalExternalChoice struct {
	currentIndex int
	InChan       chan byte
	BoxesChans   []chan byte
}

func (choice *LocalExternalChoice) WriteWithChoice() {
	if len(choice.BoxesChans) == 0 {
		panic("LocalExternalChoice requires at least one worker channel")
	}
	for _, workerInput := range choice.BoxesChans {
		if workerInput == nil {
			panic("LocalExternalChoice worker channel must not be nil")
		}
	}

	go func() {
		defer func() {
			for _, workerInput := range choice.BoxesChans {
				close(workerInput)
			}
		}()

		for msg := range choice.InChan {
			Debug("LocalExternalChoice received payload", "msg", msg)
			choice.BoxesChans[choice.currentIndex] <- msg
			choice.currentIndex = (choice.currentIndex + 1) % len(choice.BoxesChans)
		}
	}()
}

type LocalSplitBuffer struct {
	InChan  chan byte
	OutChan chan byte
}

func (smoother *LocalSplitBuffer) Interleave() {
	go func() {
		defer close(smoother.OutChan)
		for msg := range smoother.InChan {
			Debug("LocalSplitBuffer received payload", "msg", msg)
			smoother.OutChan <- msg
		}
	}()
}

type Node struct {
	ID     int
	InBox  *BoxB
	InChan chan byte
}

type LocalSplitNode struct {
	Nodes  []Node
	Buffer LocalSplitBuffer
}

func (node *LocalSplitNode) Process() {
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

type PrefixGenerator struct {
	Func    PrefixGeneratorFunc
	OutChan chan byte
}

func (generator *PrefixGenerator) Start() {
	go generator.Func(generator.OutChan)
}

type Sink struct {
	Func   SinkFunc
	InChan chan byte
}

func (sink *Sink) Consume() {
	for msg := range sink.InChan {
		sink.Func(msg)
	}
}
