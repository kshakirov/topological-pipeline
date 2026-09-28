package main

func ParallelComposition(config TcpSplitNodeConfig, f BoxFuncB) TcpSplitNode{
	//here will be a function provided and number of nodes as parallelism
	box1 := NewBoxB("b1", f)
	box2 := NewBoxB("b2", f)
	node1 := Node{ID: 1, InBox: box1, InChan: make(chan byte)}
	node2 := Node{ID: 2, InBox: box2, InChan: make(chan byte)}
	smoother := TcpSplitBuffer{InChan: make(chan byte), Config: config.BufferConfig}
	dispatcher := TcpExternalChoice{
		InChan:     make(chan byte),
		BoxesChans: []chan byte{node1.InChan, node2.InChan},
		Config:     config.ChoiceConfig,
	}

	return  TcpSplitNode{Nodes: []Node{node1, node2}, Buffer: smoother, Choice: dispatcher}
	//can be later returned to start somewhere else
}

func (pc *TcpSplitNode) Start() {
	pc.Buffer.Interleave()
	pc.Process()
	pc.Choice.WriteWithChoice()
}
