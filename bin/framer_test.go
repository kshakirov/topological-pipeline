package main

import (
	"testing"
)

func TestFramer(t *testing.T) {

	frame1 := []byte{
		0x53, 0x54, 0x00, 0x1C, // ST, Length=28
		0x04, 0x01, 0x00, // DATA, Version, Flags
		0, 0, 0, 0, 0, 0, 0, 1, // RootID
		0, 0, 0, 0, 0, 0, 0, 1, // MessageID
		0, 0, 0, 0, 0, 0, 0, 0, // ParentID
		0xAA,
	}

	// frame2 := []byte{
	// 	0x53, 0x54, 0x00, 0x1C,
	// 	0x04, 0x01, 0x00,
	// 	0, 0, 0, 0, 0, 0, 0, 1,
	// 	0, 0, 0, 0, 0, 0, 0, 2,
	// 	0, 0, 0, 0, 0, 0, 0, 0,
	// 	0xBB,
	// }

	part_1 := frame1[0:10]
	//	part_2 := frame1[10:]
	msgs, tail := Framer(part_1)
	if msgs == nil && tail == nil {
		t.Fatalf("msgs and tail  cannot be nil together \n")
	}

}
