package main

import (
	"sync"
	_ "time"
)

//оператор внешнего выбора -External Choice

type BoxFuncB func(in byte) byte

// Box — инфраструктурный контейнер (Универсальный Узел)

type BoxB struct {
	ID        string
	UserFuncB BoxFuncB
}

type PrefixGeneratorFunc func(OutChan chan byte)
type SinkFunc func(byte)

func NewBoxB(id string, fn BoxFuncB) *BoxB {
	return &BoxB{
		ID:        id,
		UserFuncB: fn,
	}
}

type LocalExternalChoice struct {
	currentIndex int
	InChan       chan byte
	BoxesChans   []chan byte
}

func (lsd *LocalExternalChoice) WriteWithChoice() {
	go func() {
		defer func() {
			for _, c := range lsd.BoxesChans {
				close(c)
			}
		}()
		for msg := range lsd.InChan {
			Debug("LocalExternalChoimce: Recieved from PrefixGenerator Payload: ", "msg", msg)
			lsd.BoxesChans[lsd.currentIndex] <- msg
			lsd.currentIndex = (lsd.currentIndex + 1) % len(lsd.BoxesChans)

		}
	}()
	//take current index пиши по модулю
}

type LocalSplitBuffer struct {
	InChan  chan byte
	OutChan chan byte
	//InChan chan byte

}

func (lsb *LocalSplitBuffer) Interleave() {
	// push all to the buffer the buffer is somehow connected to Wire
	go func() {
		defer close(lsb.OutChan)
		for msg := range lsb.InChan {
			Debug("LocalSplitBuffer: received ", "msg", msg)
			//lsb.buffer = append(lsb.buffer, msg)
			lsb.OutChan <- msg
		}
	}()
}

type Node struct {
	Id    int
	InBox *BoxB
	//OutBox *Box
	InChan  chan byte // <--- ОДИН ЕДИНСТВЕННЫЙ РАЗЪЕМ!
	OutChan chan byte
}

//g

type LocalSplitNode struct {
	ExternalChoice LocalExternalChoice
	Nodes          []Node
	Buffer         LocalSplitBuffer
}

// func (lsp *LocalSplitNode) Init(){
// 	for _,n:= range lsp.Nodes {
// 		n.OutChan = lsp.Buffer.InChan
// 	}
// }

// func (lsb * LocalSplitBuffer) Dump(b byte){
// 	lsb.buffer = append( lsb.buffer, b)
// }

func (lsp *LocalSplitNode) Process() {
	go func() {
		//Debug("Inside LocalSplitNode\n")
		//временно пока последовательно перебираем каналы
		var wg sync.WaitGroup
		for _, b := range lsp.Nodes {
			wg.Add(1)
			go func(Node) {
				defer wg.Done()
				for msg := range b.InChan {
					// Твоя рабочая двухтактная логика:
					//				Debug("msg rec\n")
					res := b.InBox.UserFuncB(msg)
					Debug("Box"," Processing Byte ", "id", b.Id, "result", res)
					lsp.Buffer.InChan <- res
					//lw.OutChan <- res
				}
			}(b)
		}
		wg.Wait()
		close(lsp.Buffer.InChan)
	}()
}

type PrefixGenerator struct {
	Func    PrefixGeneratorFunc
	OutChan chan byte
}
type Sink struct {
	Func   SinkFunc
	InChan chan byte
}

func (s *Sink) Consume() {
	//go func(){
	for msg := range s.InChan {
		s.Func(msg)

	}
	//}()
}

func (pg *PrefixGenerator) Start() {

	go pg.Func(pg.OutChan)
}
