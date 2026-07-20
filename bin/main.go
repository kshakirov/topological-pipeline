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

// BoxFunc — чистая пользовательская функция. 
// Принимает объект, возвращает трансформированный объект.






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
