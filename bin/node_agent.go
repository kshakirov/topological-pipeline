package main

import (
	_ "fmt"
	"net"
	"net/netip"
)

type NodeAgent struct {
	TopologyId         int
	NodeId             int
	DataAddress        netip.AddrPort
	CoordinatorAddress netip.AddrPort
	CoordinatorAlias string 
}

func (na *NodeAgent) Connect() {
	registration_data := make([]byte, 0, 0) //later there will be a message
	coordinator_data := make([]byte, 0, 1024)
	conn, err := net.Dial("tcp", na.CoordinatorAlias)
	if err != nil {

		Debug("Cannot Connect to Coordinator  ", "address", na.CoordinatorAddress)
		//here must be some logic to retry
	}else{
		wb, err := conn.Write(registration_data) //write can be in batches
		if err != nil {
			Debug("Cannot Write to Coordinator ", "address", na.CoordinatorAddress)
		}
		Debug("Written to Coordinator ", "bytes", wb)
		rb, err := conn.Read(coordinator_data) //read can be in batches 
		if err != nil {

			Debug("Cannot Read From Coordinator ")
		}
		Debug("Read From Coordinator ", "bytes", rb)
	}
}
