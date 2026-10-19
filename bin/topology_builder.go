package main

import (
	"log"
	"net"
)

//this one will build the topology but for the time being the only responsability is to submit topology to OTP BEAM

type Builder struct {
	Id int
}

func (builder *Builder) BuildTopology() {

}

func (builder *Builder) RegisterTopology(hostnameAlias string,  topologyMap []byte) {
	conn, err := net.Dial("tcp", hostnameAlias)
	if err != nil {
		log.Fatal("Can't register the topology, dying ...")
	}
	Debug("Connected to the orchestrator, sending the topology")
	msgLength := len(topologyMap)
	writtenBytes, err := conn.Write(topologyMap)
	for {
		if writtenBytes < msgLength {
			n, _ := conn.Write(topologyMap[writtenBytes:])
			Debug("Written bytes to Coordinator ", "bytes: ", n)
			writtenBytes += n
		}else{
			Debug("Written all  bytes to Coordinator ", "bytes: ", writtenBytes)
			break
		}
	}
	response := make([]byte, 1)
	_,err =  conn.Read(response)
	if err != nil{
		log.Fatal("The topology is not registered")
	}
	Debug("Recieved from the Orchestrator", "resonse", response)
	if response[0] != 0{
		log.Fatal("The orchestrator reject to register topology")
	}

}
