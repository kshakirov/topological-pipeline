package main

import (
	_ "fmt"
	"net"
	"net/netip"
	"slices"
)

type NodeAgent struct {
	TopologyId         int
	NodeId             int
	OutChan chan byte
	InChan chan byte
	CoordinatorAddress netip.AddrPort
	CoordinatorAlias string 
}

func revcieve_loop(conn net.Conn, outChan chan byte){
 
	dataBuffer:=make([]byte,1024)
	var tail []byte
	tail = nil
	defer close(outChan)
	for  {
		rb,err := conn.Read(dataBuffer)
		if err!= nil {
			Debug("NodeAgent: Can't read message from Coordinator ", "err", err.Error())
			break
		}
		//here checking what kind of data it is if aux data respond  else pass to Outchant
		var frames [][]byte
		var new_tail[] byte
		
		if tail!= nil {
			frames, new_tail = Framer(slices.Concat(tail, dataBuffer[0:rb]))
		}else {
			frames, new_tail = Framer(dataBuffer[0:rb])
		}

		if frames != nil {
			for _, frame := range frames {
				frame_length:= len(frame)
				for i:=0; i < frame_length; i++ {
					outChan <- frame[i]
				}
			}
		}
		if new_tail != nil {
			tail = make([]byte, len(new_tail))
			copy(tail, new_tail) 
			//do some writes info coordinator
		}else {
			tail = nil
		}
	}
}

func send_loop(conn net.Conn, inChan chan byte){
	for msg:=  range inChan {
		_,err:=conn.Write([]byte{msg})
		if err != nil {
			Debug("NodeAgent: Can't writer message to Coordinator ", "error", err.Error())
		}
		//here we just write all we have from smoother
	}
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
		//if all okayt the connection is established we can go into the background
		go revcieve_loop(conn, na.OutChan)
		go send_loop(conn, na.InChan)
		
	}
}
