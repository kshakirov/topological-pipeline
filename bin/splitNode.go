package main

import (
	"log"
_	"time"
)

//оператор внешнего выбора -External Choice


type BoxFuncB func(in byte) byte

// Box — инфраструктурный контейнер (Универсальный Узел)

type BoxB struct {
	ID          string
	UserFuncB    BoxFuncB
	
}

type PrefixGeneratorFunc func(OutChan chan byte)


func NewBoxB(id string, fn BoxFuncB) *BoxB {
	return &BoxB{
		ID:          id,
		UserFuncB:    fn,

	}
}


type LocalExternalChoice struct{
	currentIndex int
	indices []int
	channelsNumber int
	InChan chan byte
	BoxesChans []chan byte
	//Wire который надо передать сорсу
	
}

func (lsd *LocalExternalChoice) WriteWithChoice (){
	go func(){
		for msg := range lsd.InChan {
			log.Printf("LocalExternalChoime: Recieved from PrefixGenerator Payload: [%d] \n", msg)
			for i := range lsd.indices{
				lsd.BoxesChans[i] <- msg
				lsd.currentIndex = (lsd.currentIndex + 1) % len(lsd.BoxesChans)
			}

			
		}
	}()
	//take current index пиши по модулю 
}

type LocalSplitBuffer struct{
	buffer []byte
	//InChan chan byte
	
}

func (lsp *LocalSplitBuffer) Interleave(out byte){
	// push all to the buffer the buffer is somehow connected to Wire
}

type Node struct {
	Id     int
	InBox  *BoxB
	//OutBox *Box
	InChan   chan byte // <--- ОДИН ЕДИНСТВЕННЫЙ РАЗЪЕМ!
	OutChan chan byte
	
}
//g


type LocalSplitNode struct{
	ExternalChoice LocalExternalChoice
	Nodes []Node
	Buffer LocalSplitBuffer
}

//func (lsn *LocalSplitNode) Writer (Set) 


func (lsb * LocalSplitBuffer) Dump(b byte){
	lsb.buffer = append( lsb.buffer, b)
}



func (lsp *LocalSplitNode) Process() {
	go func() {
		//log.Printf("Inside LocalSplitNode\n")
		//временно пока последовательно перебираем каналы
		for _,b := range lsp.Nodes {
			go func(Node){
				for msg := range b.InChan {
					// Твоя рабочая двухтактная логика:
					//				log.Printf("msg rec\n")
					res := b.InBox.UserFuncB(msg)
					log.Printf("b id is %d %v\n",b.Id, res)
					lsp.Buffer.Dump(res)
					//lw.OutChan <- res
				}
			}(b)
		}
	}()
}


type PrefixGenerator struct{
	Func  PrefixGeneratorFunc
	OutChan chan byte 
}


func (pg *PrefixGenerator)Start(b byte){
	
	go pg.Func(pg.OutChan)
}

