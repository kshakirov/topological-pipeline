package main

import (
	"context"
	"os"
	"os/signal"
)

func main() {
	Debug("Starting Go Storm")

	box1 := NewBoxB("b1", testComputeB)
	box2 := NewBoxB("b2", testComputeB)
	node1 := Node{ID: 1, InBox: box1, InChan: make(chan byte)}
	node2 := Node{ID: 2, InBox: box2, InChan: make(chan byte)}
	smoother := LocalSplitBuffer{InChan: make(chan byte), OutChan: make(chan byte)}
	dispatcher := LocalExternalChoice{
		InChan:     make(chan byte),
		BoxesChans: []chan byte{node1.InChan, node2.InChan},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	generator := PrefixGenerator{Func: testPrefixGeneratorLoop(ctx), OutChan: dispatcher.InChan}
	parallelNode := LocalSplitNode{Nodes: []Node{node1, node2}, Buffer: smoother}
	sink := Sink{InChan: smoother.OutChan, Func: testSinkFunc}

	smoother.Interleave()
	parallelNode.Process()
	dispatcher.WriteWithChoice()
	generator.Start()
	sink.Consume()

	Debug("Go Storm stopped gracefully")
}
