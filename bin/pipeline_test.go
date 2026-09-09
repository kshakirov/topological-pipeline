package main

import (
	"context"
	"testing"
	"time"
)

func TestFinitePipelineDrainsToSink(t *testing.T) {
	const inputCount = 128

	box1 := NewBoxB("b1", testComputeB)
	box2 := NewBoxB("b2", testComputeB)
	node1 := Node{ID: 1, InBox: box1, InChan: make(chan byte)}
	node2 := Node{ID: 2, InBox: box2, InChan: make(chan byte)}
	smoother := LocalSplitBuffer{InChan: make(chan byte), OutChan: make(chan byte)}
	dispatcher := LocalExternalChoice{
		InChan:     make(chan byte),
		BoxesChans: []chan byte{node1.InChan, node2.InChan},
	}
	parallelNode := LocalSplitNode{
		Nodes:  []Node{node1, node2},
		Buffer: smoother,
	}

	smoother.Interleave()
	parallelNode.Process()
	dispatcher.WriteWithChoice()

	go func() {
		defer close(dispatcher.InChan)
		for i := 0; i < inputCount; i++ {
			dispatcher.InChan <- byte(i)
		}
	}()

	var total, zeros, ones int
	for result := range smoother.OutChan {
		total++
		switch result {
		case 0:
			zeros++
		case 1:
			ones++
		default:
			t.Fatalf("unexpected result %d", result)
		}
	}

	if total != inputCount || zeros != inputCount/2 || ones != inputCount/2 {
		t.Fatalf("got total=%d zeros=%d ones=%d; want total=128 zeros=64 ones=64", total, zeros, ones)
	}
}

func TestCancelledGeneratorClosesOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	out := make(chan byte)
	done := make(chan struct{})
	go func() {
		testPrefixGeneratorLoop(ctx)(out)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("generator did not stop after cancellation")
	}

	if _, open := <-out; open {
		t.Fatal("generator output is still open after cancellation")
	}
}
