package main

func main() {
	Debug("Starting Go Storm")
	generator := PrefixGeneratorTCP{Func: testPrefixGeneratorLoopTCP("127.0.0.1", 9090)}
	consumer := SinkTCP{Func: testSinkGeneratorLoopTCP("127.0.0.1",9091)}

	
	ParallelComposition(TcpSplitNodeConfig{ TcpExternalChoiceConfig{"127.0.0.1",9090},
		TcpSpliterBuffferConfig{"127.0.0.1",9091}, 2})
	
	// box1 := NewBoxB("b1", testComputeB)
	// box2 := NewBoxB("b2", testComputeB)
	// node1 := Node{ID: 1, InBox: box1, InChan: make(chan byte)}
	// node2 := Node{ID: 2, InBox: box2, InChan: make(chan byte)}
	// smoother := TcpSplitBuffer{InChan: make(chan byte), Config: TcpSpliterBuffferConfig{"127.0.0.1",9091}}
	// dispatcher := TcpExternalChoice{
	// 	InChan:     make(chan byte),
	// 	BoxesChans: []chan byte{node1.InChan, node2.InChan},
	// 	Config: TcpExternalChoiceConfig{"127.0.0.1",9090},
	// }
	// parallelNode := TcpSplitNode{Nodes: []Node{node1, node2}, Buffer: smoother, Choice: dispatcher}
	// smoother.Interleave()
	// parallelNode.Process()
	// dispatcher.WriteWithChoice()

	generator.Start()
	consumer.Consume()
	Debug("Go Storm stopped gracefully")
}
