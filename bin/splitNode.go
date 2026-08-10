package main

import (
	"log"
	"time"
)

//оператор внешнего выбора -External Choice


type BoxFuncB func(in byte) byte

// Box — инфраструктурный контейнер (Универсальный Узел)

type BoxB struct {
	ID          string
	UserFuncB    BoxFuncB
	
}




func NewBoxB(id string, fn BoxFuncB) *BoxB {
	return &BoxB{
		ID:          id,
		UserFuncB:    fn,

	}
}


type LocalSplitDispatcher struct{
	currentIndex int
	indices []int
	channelsNumber int
	InChan chan byte
	BoxesChans []chan byte
	//Wire который надо передать сорсу
	
}

func (lsd *LocalSplitDispatcher) WriteWithChoice (){
	for i := range lsd.indices{
		time.Sleep(time.Second * 5)
		lsd.currentIndex = i
		b := byte(i)
		lsd.BoxesChans[i] <- b
	}
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
	Dispatcher LocalSplitDispatcher
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
			for msg := range b.OutChan {
				// Твоя рабочая двухтактная логика:
				//				log.Printf("msg rec\n")
				res := b.InBox.UserFuncB(msg)
				log.Printf("%v\n", res)
				lsp.Buffer.Dump(res)
				//lw.OutChan <- res
			}
		}
	}()
}
