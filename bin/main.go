package main

import (

	"log"
  	_ "sync"
	"time"
)






func main() {

	
	box := NewBox("QuantumProcessor", heavyCompute)
	c_box_1:=NewBox("testOut1", testCompute)
	c_box_2:=NewBox("testOut2", testCompute)
	sourceNode:=SourceNode{OutBox:box}
	computeNode := ComputeNode{InBox:c_box_1 , OutBox: c_box_2}
	log.Printf("%v\n", computeNode)
	computeNode.Prep()
	computeNode.Wire.WireIn(computeNode.InBox,computeNode.OutBox)
	sourceNode.AddChannel(&computeNode.OutChan)
	
	sourceNode.Start(23)
	time.Sleep(time.Second * 2)
}

