package main

import (
_"log"


)

func (s *ComputeNode) AddInOut(i *Box, o *Box){
	s.OutBox = o
	s.InBox = i
}

func (s *ComputeNode) Prep(){
	lw:=&LocalWire{Id:1}
	linw:=&LocalInWire{InChan: make(chan Set, 10)}
 	loutw:=&LocalOutWire{OutChan: make(chan Set, 10)}
	lw.InWire = &LocalInWire{}
	lw.OutWire = loutw
	lw.InWire = linw
	
	s.Wire = lw
	lw.InWire.Init(10)
	lw.OutWire.Init(10)
		//s.InChan = make(chan Set, 10)
		//s.OutChan = make(chan Set, 10)
	//	linw.Init(10)
	//lw.OutChan.Init(10)
	//loutw.Init(10)
	//s.InChan.Init(10)
	
}



// func (s *ComputeNode) WireIn(){
// 	go func() {
// 		for msg := range s.InChan {
// 			log.Printf("ComputeNode[%d]: The msg arrived\n",s.Id );
// 			res := s.InBox.UserFunc(msg)
// 			log.Printf("InBox [%s] applied to the msg, result is [%v]\n", s.InBox.ID, res)
// 			res = s.OutBox.UserFunc(res)
// 			log.Printf("OutBox [%s] applied to the msg result is [%v]\n", s.OutBox.ID, res)
// 			s.OutChan <- res
// 		}
// 	}()

// }
