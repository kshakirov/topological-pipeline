package main

import "errors"



type Message struct {
	RootId [16] byte
	MessageId [16]byte
	ParentId [16]byte
	Flags byte
	Payload []byte
}

func Decode(payload []byte ) (data Message, e error ){
	var msg Message
	if len(payload ) >= 51 && payload [0] == 4 && payload[1]==1{
		copy(msg.RootId[:], payload[3:])
		copy(msg.MessageId[:], payload[19:])
		copy(msg.ParentId[:], payload[35:])
		msg.Payload= make([]byte, len(payload) - 51)
		copy(msg.Payload[:], payload[51:])
		
		return msg,nil 
	}else{
		
	}
	return msg, errors.New("Msg does not comply with the format") ;
}
	
