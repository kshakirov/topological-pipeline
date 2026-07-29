package main

import (
	"log"
	_ "log"
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
	//	InChan *chan Set
	Wire Wire

}


type ComputeNode struct {
    Id     int
    InBox  *Box
    OutBox *Box
    Wire   Wire // <--- ОДИН ЕДИНСТВЕННЫЙ РАЗЪЕМ!
}


type Tuple struct{
	Fst any
	Snd any
}

type LocalTuple struct {
	Fst Set
	Snd Set
}

type InWire interface{
	 Init(any)
}



type OutWire interface{
	Write(Tuple, bool)
	Init(any)
}



type Wire interface {
    Read() Set               // Забрать квант данных (для такта IN)
    Write(msg Set)           // Протолкнуть квант данных (для такта OUT)
    WireIn(inBox Box, outBox Box) // Запустить внутреннюю горутину связи
}






type LocalWire struct {
    Id      int
    InChan  chan Set
    OutChan chan Set
}

func (lw *LocalWire) Read() Set {
    return <-lw.InChan
}

func (lw *LocalWire) Write(msg Set) {
	log.Printf("Write: %v\n",msg)
	lw.OutChan <- msg
}

func (lw *LocalWire) WireIn(inBox Box, outBox Box) {
	go func() {
		log.Printf("inside localwire:\n")
        for msg := range lw.OutChan {
		// Твоя рабочая двухтактная логика:
		log.Printf("msg rec\n")
            res := inBox.UserFunc(msg)
            res = outBox.UserFunc(res)
		//lw.OutChan <- res
        }
    }()
}
