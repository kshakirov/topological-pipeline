package main

import (

	"log"

	_ "gostorm.org/go_storm/lib/bolt"
	_ "gostorm.org/go_storm/lib/spout"
	_ "gostorm.org/go_storm/lib/tuple"

	_ "sync"
	"time"
)

// Data — универсальный объект нашей категории (пока просто interface{})
type Set interface{
}

// BoxFunc — чистая пользовательская функция. 
// Принимает объект, возвращает трансформированный объект.
type BoxFunc func(in Set) Set

// Box — инфраструктурный контейнер (Универсальный Узел)

type Box struct {
	ID          string
	UserFunc    BoxFunc
	
	// может будет что то еще здесь, наверняка
	//мультиплекоср или демультиплпксор будет содержать всю магию
}
type MuxNode struct {
	InBoxes []Box
	OutBox *Box
}



type ComputeNode struct{
	Id int
	OutBox *Box
	InBox  *Box
	InChan chan Set
	OutChan chan Set
	
}

type SourceNode struct{
	OutBox *Box
	InChan *chan Set

	
}


func (s *SourceNode) AddOut(o *Box){
	s.OutBox = o

}


func (s *SourceNode) Start(input Set){
	log.Printf("StartNode: Input = %v\n", s)
	res:= s.OutBox.UserFunc(input)
	log.Printf("StartNode: Output= %v \n", res)
	*(s.InChan) <- res


}



func (s *ComputeNode) AddInOut(i *Box, o *Box){
	s.OutBox = o
	s.InBox = i
}

func (s *ComputeNode) Prep(){
	s.InChan = make(chan Set, 10)
	s.OutChan = make(chan Set, 10)
}


func (s *ComputeNode) WireIn(){
	go func() {
		for msg := range s.InChan {
			log.Printf("ComputeNode[%d]: The msg arrived\n",s.Id );
			res := s.InBox.UserFunc(msg)
			log.Printf("InBox [%s] applied to the msg, result is [%v]\n", s.InBox.ID, res)
			res = s.OutBox.UserFunc(res)
			log.Printf("OutBox [%s] applied to the msg result is [%v]\n", s.OutBox.ID, res)
			s.OutChan <- res
		}
	}()

}

type Vertex interface {
	WireIn(ch<- chan Set)
}

type Graph struct {
	Vertices []Vertex
}


func (g *Graph)Run(){
	
}

type DemuxNode struct {
	OutBoxes []Box
	InBox *Box
}


func (d *DemuxNode) AddIn(box *Box){
	d.InBox = box
}

func (d *DemuxNode) AddOut(box Box){
	d.OutBoxes = append(d.OutBoxes, box)
}


func (d *MuxNode) AddIn(box *Box){
	d.OutBox = box
}

func (d *MuxNode) addOut(box Box){
	d.InBoxes = append(d.InBoxes, box)
}



func NewBox(id string, fn BoxFunc) *Box {
	return &Box{
		ID:          id,
		UserFunc:    fn,

	}
}

// Start запускает схемотехнику Ящика
func (b *Box) Start(msg Set ) Set{
	
	res := b.UserFunc(msg)
	return res
	

}

func (s *SourceNode) AddChannel(c *chan Set){
	s.InChan = c

}



func wireStartNode(s *SourceNode,c *ComputeNode){
	s.AddChannel(&c.InChan)
	
}


func main() {

	
	box := NewBox("QuantumProcessor", heavyCompute)
	c_box_1:=NewBox("testOut1", testCompute)
	c_box_2:=NewBox("testOut2", testCompute)
	sourceNode:=SourceNode{OutBox:box}
	computeNode := ComputeNode{InBox:c_box_1 , OutBox: c_box_2}
	
	computeNode.Prep()
	computeNode.WireIn()
	sourceNode.AddChannel(&computeNode.InChan)
	
	sourceNode.Start(23)
	time.Sleep(time.Second * 2)
}

func heavyCompute(in Set) Set {
	val := in.(int)
	log.Printf("Func heavyCompute input=%d\n", val);
	// Искусственный разброс по времени, чтобы проверить FIFO на выходе
	if val%2 == 0 {
		time.Sleep(50 * time.Millisecond) // Четные — тугодумы
	} else {
		time.Sleep(10 * time.Millisecond) // Нечетные — шустрики
	}
	return val

}


func testCompute(i Set) Set{
	log.Printf("Func TestCompute input= %d\n", i)
	val := i.(int)
	return val * 20
}
