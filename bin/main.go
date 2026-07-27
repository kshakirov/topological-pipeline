package main

import (

	"log"
  	_ "sync"
	"time"
)






func main() {

	
	box := NewBox("QuantumProcessor", heavyCompute)
	inBox:=NewBox("testOut1", testCompute)
	outBox:=NewBox("testOut2", testCompute)
	sourceNode:=SourceNode{OutBox:box}
	//	computeNode := ComputeNode{InBox:c_box_1 , OutBox: c_box_2}
	myWire := &LocalWire{Id: 1, InChan: make(chan Set,10), OutChan: make(chan Set, 10)}
	
	computeNode := ComputeNode{ Id: 0, InBox: inBox, OutBox: outBox, Wire: myWire }

	log.Printf("%v\n", computeNode)
	computeNode.Prep()
	//computeNode.Wire.WireIn(computeNode.InBox,computeNode.OutBox)
	sourceNode.AddChannel(computeNode.Wire)

	sourceNode.Start(23)
	time.Sleep(time.Second * 2)
}

