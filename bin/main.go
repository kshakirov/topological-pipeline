package main

import (
	_ "sync"
	_ "time"
)

func main() {

	box_1 := NewBoxB("b1", testComputeB)
	box_2 := NewBoxB("b2", testComputeB)
	node_1 := Node{Id: 1, InBox: box_1, InChan: make(chan byte), OutChan: make(chan byte)}
	node_2 := Node{Id: 2, InBox: box_2, InChan: make(chan byte), OutChan: make(chan byte)}
	buffer := LocalSplitBuffer{InChan: make(chan byte), OutChan: make(chan byte)}

	dispatcher := LocalExternalChoice{currentIndex: 0, InChan: make(chan byte), BoxesChans: []chan byte{node_1.InChan, node_2.InChan}}
	pfg := PrefixGenerator{Func: testPrefixGeneratorLoop(), OutChan: dispatcher.InChan}
	lsp := LocalSplitNode{ExternalChoice: dispatcher, Nodes: []Node{node_1, node_2}, Buffer: buffer}
	sink := Sink{InChan: buffer.OutChan, Func: testSinkFunc}

	buffer.Interleave()
	lsp.Process()

	dispatcher.WriteWithChoice()
	pfg.Start()
	sink.Consume()

}
