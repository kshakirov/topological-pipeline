package main

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestReceiveLoopFragments(t *testing.T) {
	reader, writer := net.Pipe()
	defer reader.Close()
	defer writer.Close()
	out := make(chan byte)
	go revcieve_loop(reader, out)
	go func() {
		defer writer.Close()
		for _, part := range [][]byte{
			{0x53, 0x54, 0, 2, 4},
			{0xAA, 0x53, 0x54, 0, 3, 4},
			{0xBB, 0x0A},
		} {
			if _, err := writer.Write(part); err != nil {
				return
			}
		}
	}()

	var got []byte
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case b, ok := <-out:
			if !ok {
				if !bytes.Equal(got, []byte{4, 0xAA, 4, 0xBb, 0x0A}) {
					t.Fatalf("got % X", got)
				}
				return
			}
			got = append(got, b)
		case <-timer.C:
			t.Fatal("receive loop did not finish")
		}
	}
}
