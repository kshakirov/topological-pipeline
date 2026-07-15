package main

import (
	"fmt"
	_ "gostorm.org/go_storm/lib/spout"
	_ "gostorm.org/go_storm/lib/bolt"
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
	OutBox *Box
	InBox  *Box
	InChan chan Set
	
}

type SourceNode struct{
	OutBox *Box
	InChan *chan Set
	
}


func (s *SourceNode) AddOut(o *Box){
	s.OutBox = o

}


func (s *SourceNode) Start(input Set){
	fmt.Printf("%v\n", s)
	*(s.InChan) <- s.OutBox.UserFunc(input)

}



func (s *ComputeNode) AddInOut(i *Box, o *Box){
	s.OutBox = o
	s.InBox = i
}

func (s *ComputeNode) Prep(){
	s.InChan = make(chan Set, 10)
}


func (s *ComputeNode) WireIn(){
	go func() {
		for msg := range s.InChan {
			fmt.Printf("here we go");
			res := s.OutBox.UserFunc(msg)
			s.InChan <- res
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
	c_box:=NewBox("test", testCompute)
	sourceNode:=SourceNode{OutBox:box}
	computeNode := ComputeNode{InBox: box, OutBox: c_box}
	computeNode.WireIn()
	computeNode.Prep()
	sourceNode.AddChannel(&computeNode.InChan)
	sourceNode.Start(1)
	// for i := 1; i <= 6; i++ {
	// 	fmt.Printf("📥 Подано на вход: %d\n", i)
	// 	res:= box.Start(i)
	// 	fmt.Printf("Получено со входа %s\n",res)
	// }

	// fmt.Println("🏁 Конвейер полностью остановлен.")
}

func heavyCompute(in Set) Set {
	val := in.(int)
	fmt.Printf("%d in heavy \n", val);
	// Искусственный разброс по времени, чтобы проверить FIFO на выходе
	if val%2 == 0 {
		time.Sleep(50 * time.Millisecond) // Четные — тугодумы
	} else {
		time.Sleep(10 * time.Millisecond) // Нечетные — шустрики
	}
	return fmt.Sprintf("Result(%d)", val*10)
}


func testCompute(i Set) Set{
	fmt.Printf("test %d\n", i)
	val := i.(int)
	return val * 100
}
