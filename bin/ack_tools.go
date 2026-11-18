package main

import "github.com/google/uuid"

type CompleteMessage struct {
	rootId      uuid.UUID
	completedId uuid.UUID
	children    []uuid.UUID
}

func Begin(rootId uuid.UUID, msgId uuid.UUID) []byte {
	msg := make([]byte, 33)
	msg[0] = 7
	copy(msg[1:17], rootId[:])
	copy(msg[17:], msgId[:])
	return msg
}

func Complete(completeMsg CompleteMessage) []byte {
	length := 35 + len(completeMsg.children)*16
	count := len(completeMsg.children)
	msg := make([]byte, length)
	msg[0] = 8
	copy(msg[1:17], completeMsg.rootId[:])
	copy(msg[17:], completeMsg.completedId[:])
	msg[33] = byte(count >> 8)
	msg[34] = byte(count)
	for i, u := range completeMsg.children {
		since := 35 + i*16
		copy(msg[since:], u[:])
	}
	return msg
}
