package main

import (

	"log"
  	_ "sync"
_	"time"
)






func main() {

	
	// box := NewBox("QuantumProcessor", heavyCompute)
	// inBox:=NewBox("testOut1", testCompute)
	// outBox:=NewBox("testOut2", testCompute)
	// sourceNode:=SourceNode{OutBox:box}
	// //	computeNode := ComputeNode{InBox:c_box_1 , OutBox: c_box_2}
	// myWire := &LocalWire{Id: 1, InChan: make(chan Set,10), OutChan: make(chan Set, 10)}
	
	// computeNode := ComputeNode{ Id: 0, InBox: inBox, OutBox: outBox, Wire: myWire }

	// log.Printf("%v\n", computeNode)
	// computeNode.Prep()
	// //computeNode.Wire.WireIn(computeNode.InBox,computeNode.OutBox)
	// sourceNode.AddChannel(computeNode.Wire)

	// sourceNode.Start(23)
	// time.Sleep(time.Second * 2)
	box_1 := NewBoxB("b1", testComputeB)
	box_2 := NewBoxB("b2", testComputeB)
	node_1 := Node{Id: 1, InBox: box_1, InChan: make(chan byte), OutChan: make(chan byte)}
	node_2 := Node{Id: 2, InBox: box_2, InChan: make(chan byte), OutChan: make(chan byte)}
	buffer:= LocalSplitBuffer{buffer: make([]byte, 256)}

	dispatcher := LocalSplitDispatcher{currentIndex: 0, indices: []int{1,2}, channelsNumber: 2, InChan: make(chan byte), BoxesChans: []chan byte{node_1.InChan, node_2.InChan}}
	lsp:= LocalSplitNode{Dispatcher: dispatcher, Nodes: []Node{node_1, node_2}, Buffer: buffer}
	dispatcher.WriteWithChoice()
	lsp.Process()
	log.Printf("%v, %v, %v", dispatcher, box_1, box_2)
	

}

