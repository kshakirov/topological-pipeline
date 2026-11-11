package main

func Encoder(msg Message)([]byte){
	length := 51 + len(msg.Payload)
	payload := make([]byte, length)
	return payload
}
