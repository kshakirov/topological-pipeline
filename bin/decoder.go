package main

import (
	"errors"

	"github.com/google/uuid"
)

type Message struct {
	RootId    uuid.UUID
	MessageId uuid.UUID
	ParentId  uuid.UUID
	Flags     byte
	Payload   []byte
}

func Decoder(payload []byte) (data Message, e error) {
	var msg Message
	if len(payload) >= 51 && payload[0] == 4 && payload[1] == 1 {
		id, err := uuid.FromBytes(payload[3:19])
		if err == nil {
			msg.RootId = id
		}
		msg.MessageId, err = uuid.FromBytes(payload[19:35])
		msg.ParentId, err = uuid.FromBytes(payload[35:51])
		msg.Flags = payload[2]
		msg.Payload = make([]byte, len(payload)-51)
		copy(msg.Payload[:], payload[51:])

		return msg, nil
	} else {

	}
	return msg, errors.New("Msg does not comply with the format")
}
