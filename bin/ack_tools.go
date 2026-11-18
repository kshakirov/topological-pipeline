package main

import "github.com/google/uuid"

func Begin(rootId uuid.UUID, msgId uuid.UUID) ([]byte){
	msg := make([]byte, 33)
	msg[0]=7
	copy(msg[1:17], rootId[:])
	copy(msg[17:], msgId[:])
	return msg
}
