package main

import "fmt"

func ParallelComposition(config TcpSplitNodeConfig, f BoxFuncB) TcpSplitNode{
	nodes:= make([]Node,config.Parallelism)
	boxes_chans :=make([]chan byte,  config.Parallelism)
	for i:=0;i < config.Parallelism; i++{
		nodes[i]=Node{ID: i + 1, InBox: NewBoxB(fmt.Sprintf("b%d",i + 1), f), InChan: make(chan byte)}
		boxes_chans[i]= nodes[i].InChan
	}
	//here will be a function provided and number of nodes as parallelism
	//	box1 := NewBoxB("b1", f)
	//	box2 := NewBoxB("b2", f)
	//	node1 := Node{ID: 1, InBox: box1, InChan: make(chan byte)}
	//	node2 := Node{ID: 2, InBox: box2, InChan: make(chan byte)}
	smoother := TcpSplitBuffer{InChan: make(chan byte), Config: config.BufferConfig}
	dispatcher := TcpExternalChoice{
		InChan:     make(chan byte),
		BoxesChans: boxes_chans,
		Config:     config.ChoiceConfig,
	}

	return  TcpSplitNode{Nodes: nodes, Buffer: smoother, Choice: dispatcher}
	//can be later returned to start somewhere else
}

func (pc *TcpSplitNode) Start() {
	pc.Buffer.Interleave()
	pc.Process()
	pc.Choice.WriteWithChoice()
}
