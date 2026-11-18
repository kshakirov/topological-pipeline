package main

import (
	_ "fmt"
	"io"
	"net"
	"net/netip"
	"slices"
)

type NodeAgent struct {
	TopologyId         int
	NodeId             int
	OutChan            chan Message
	InChan             chan Message
	CoordinatorAddress netip.AddrPort
	CoordinatorAlias   string
}

func recieve_loop(conn net.Conn, outChan chan Message) {

	dataBuffer := make([]byte, 1024)
	var tail []byte
	tail = nil
	defer close(outChan)
	for {
		rb, err := conn.Read(dataBuffer)
		if err != nil {
			Debug("NodeAgent: Can't read message from Coordinator ", "err", err.Error())
			break
		}
		//here checking what kind of data it is if aux data respond  else pass to Outchant
		var frames [][]byte
		var new_tail []byte

		if tail != nil {
			frames, new_tail = Framer(slices.Concat(tail, dataBuffer[0:rb]))
		} else {
			frames, new_tail = Framer(dataBuffer[0:rb])
		}

		if frames != nil {
			for _, frame := range frames {

				msg, e := Decoder(frame)
				if e != nil {
					Debug("Something wrong wiith decoding", "frame", frame)
				} else {
					outChan <- msg
				}

			}
		}
		if new_tail != nil {
			tail = make([]byte, len(new_tail))
			copy(tail, new_tail)
			//do some writes info coordinator
		} else {
			tail = nil
		}
	}
}

func send_loop(conn net.Conn, inChan chan Message) {
	//

	for msg := range inChan {
		payload := PrefixFrame(Encoder(msg))
		_, err := conn.Write(payload)
		if err != nil {
			Debug("NodeAgent: Can't writer message to Coordinator ", "error", err.Error())
		}
		//here we just write all we have from smoother
	}
}

func (na *NodeAgent) Connect() {
	registration_data := []byte{
		0x53, 0x54,
		0x00, 0x05,
		0x01,
		byte(na.TopologyId >> 8), byte(na.TopologyId),
		byte(na.NodeId >> 8), byte(na.NodeId),
	}
	coordinator_data := make([]byte, 5)
	conn, err := net.Dial("tcp", na.CoordinatorAlias)
	if err != nil {

		Debug("Cannot Connect to Coordinator  ", "address", na.CoordinatorAddress)
		//here must be some logic to retry
	} else {
		wb, err := conn.Write(registration_data) //write can be in batches
		if err != nil {
			Debug("Cannot Write to Coordinator ", "address", na.CoordinatorAddress)
		}
		_ = wb
		rb, err := io.ReadFull(conn, coordinator_data)
		if err != nil {

			Debug("Cannot Read From Coordinator ")
		}
		_ = rb
		//if all okayt the connection is established we can go into the background
		go recieve_loop(conn, na.OutChan)
		go send_loop(conn, na.InChan)

	}
}
