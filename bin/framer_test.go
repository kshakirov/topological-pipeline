package main

import (
	"bytes"
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
	if msgs != nil {
		t.Fatalf("msgs must be nil \n")
	}
	if tail == nil {
		t.Fatalf("tail  must not  be nil \n")
	}

	msgs, tail = Framer(frame1)
	if msgs == nil {
		t.Fatalf("msgs must not be nil \n")
	}

}

func TestFramerMultipleFrames(t *testing.T) {
	frame1 := []byte{
		0x53, 0x54, 0x00, 0x1C,
		0x04, 0x01, 0x00,
		0, 0, 0, 0, 0, 0, 0, 1, // RootID
		0, 0, 0, 0, 0, 0, 0, 1, // MessageID
		0, 0, 0, 0, 0, 0, 0, 0, // ParentID
		0xAA,
	}
	frame2 := append([]byte(nil), frame1...)
	frame2[22] = 2 // MessageID = 2
	frame2[len(frame2)-1] = 0xBB
	input := append(append([]byte(nil), frame1...), frame2...)

	msgs, tail := Framer(input)
	if len(msgs) != 2 {
		t.Fatalf("got %d payloads; want 2", len(msgs))
	}
	for i, want := range [][]byte{frame1[4:], frame2[4:]} {
		if !bytes.Equal(msgs[i], want) {
			t.Errorf("payload %d: got % X; want % X", i, msgs[i], want)
		}
	}
	if len(tail) != 0 {
		t.Errorf("unexpected tail: % X", tail)
	}
}

func TestFramerFrameAndPartialTail(t *testing.T) {
	input := []byte{0x53, 0x54, 0, 2, 4, 0xAA, 0x53, 0x54, 0, 3, 4}

	msgs, tail := Framer(input)
	if len(msgs) != 1 || !bytes.Equal(msgs[0], []byte{4, 0xAA}) {
		t.Fatalf("unexpected payloads: % X", msgs)
	}
	if !bytes.Equal(tail, []byte{0x53, 0x54, 0, 3, 4}) {
		t.Fatalf("unexpected tail: % X", tail)
	}
}
