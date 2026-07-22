package main

import (
	"log"
)

type Set interface{

}

type BoxFunc func(in Set) Set

// Box — инфраструктурный контейнер (Универсальный Узел)

type Box struct {
	ID          string
	UserFunc    BoxFunc
	
}

type MuxNode struct {
	InBoxes []Box
	OutBox *Box
}


type SourceNode struct{
	OutBox *Box
	InChan *chan Set

}


type ComputeNode struct{
	Id int
	OutBox *Box
	InBox  *Box
	InChan chan Set
	OutChan chan Set
	
}


type Tuple struct{
	Fst any
	Snd any
}


type InWire  any


type OutWire any

type Wire interface{
	OutWire
	InWire
	Commute()	
}


type LocalInWire struct {
	InChan chan Set
	
}

type LocalOutWire struct{
	OutChan chan Set
}

type LocalWire struct {
	Id int
	InChan chan Set
	OutChan chan Set
	
}

func (lw *LocalWire)WireIn(inBox Box, outBox Box){
 	go func() {
		for msg := range lw.InChan {
			log.Printf("Wire [%d]: The msg arrived\n",lw.Id );
			res := inBox.UserFunc(msg)
			log.Printf("InBox [%s] applied to the msg, result is [%v]\n", inBox.ID, res)
			res = outBox.UserFunc(res)
			log.Printf("OutBox [%s] applied to the msg result is [%v]\n", outBox.ID, res)
			lw.OutChan <- res
		}
	}()
}
