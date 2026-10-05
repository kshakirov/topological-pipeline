package main

import (
	"fmt"
)

func ParallelComposition(config TcpSplitNodeConfig, f BoxFuncB) TcpSplitNode {
	nodes := make([]Node, config.Parallelism)
	boxes_chans := make([]chan byte, config.Parallelism)
	for i := 0; i < config.Parallelism; i++ {
		nodes[i] = Node{ID: i + 1, InBox: NewBoxB(fmt.Sprintf("b%d", i+1), f), InChan: make(chan byte)}
		boxes_chans[i] = nodes[i].InChan
	}

	smoother := TcpSplitBuffer{InChan: make(chan byte), Config: config.BufferConfig, OutChan: make (chan byte)}
	node_agent := NodeAgent{TopologyId: 1, NodeId: 1, CoordinatorAlias: "elixirCoordinator:900", OutChan: make(chan byte), InChan: smoother.OutChan}
	dispatcher := TcpExternalChoice{
		InChan:     node_agent.OutChan,
		BoxesChans: boxes_chans,
		Config:     config.ChoiceConfig,
	}
	//will be provided with the function


	return TcpSplitNode{Nodes: nodes, Buffer: smoother, Choice: dispatcher, Agent: node_agent}
	//can be later returned to start somewhere else
}

func (pc *TcpSplitNode) Start() {
	pc.Agent.Connect() //just to test
	pc.Buffer.Interleave()
	pc.Process()
	pc.Choice.WriteWithChoice()
}
