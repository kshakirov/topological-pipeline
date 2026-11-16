package main

func Encoder(msg Message) []byte {
	length := 51 + len(msg.Payload)
	encoded_msg := make([]byte, length)
	copy(encoded_msg[51:], msg.Payload[:])
	copy(encoded_msg[19:], msg.MessageId[:])
	copy(encoded_msg[35:], msg.ParentId[:])
	copy(encoded_msg[3:], msg.RootId[:])
	encoded_msg[2] = msg.Flags
	encoded_msg[0] = 0x04
	encoded_msg[1] = 0x01
	return encoded_msg
}
