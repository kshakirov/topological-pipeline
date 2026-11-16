package main

import (
	"fmt"
)

func ParallelComposition(config TcpSplitNodeConfig, f BoxFuncB) TcpSplitNode {
	nodes := make([]Node, config.Parallelism)
	boxes_chans := make([]chan Message, config.Parallelism)
	for i := 0; i < config.Parallelism; i++ {
		nodes[i] = Node{ID: i + 1, InBox: NewBoxB(fmt.Sprintf("b%d", i+1), f), InChan: make(chan Message)}
		boxes_chans[i] = nodes[i].InChan
	}

	smoother := TcpSplitBuffer{InChan: make(chan Message), Config: config.BufferConfig, OutChan: make(chan Message)}
	node_agent := NodeAgent{TopologyId: 1, NodeId: 2, CoordinatorAlias: "127.0.0.1:9300", OutChan: make(chan Message), InChan: smoother.OutChan}
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
