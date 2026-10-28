package main

import (
	"testing"
)

func TestFramer(t *testing.T) {
	topology := []byte{
		0x53, 0x54, // ST
		0x00, 0x0F, // Payload length: 15 bytes
		0x03,       // CREATE_TOPOLOGY
		0x00, 0x01, // TopologyID: 1
		0x00, 0x02, // NodeCount: 2
		0x00, 0x01, // NodeID: 1
		0x00, 0x02, // NodeID: 2
		0x00, 0x01, // EdgeCount: 1
		0x00, 0x01, // FromID: 1
		0x00, 0x02, // ToID: 2
	}
	msgs, tail := Framer(topology)
	if msgs == nil && tail == nil{
		t.Fatalf("msgs and tail  cannot be nil together \n")
	}

}
