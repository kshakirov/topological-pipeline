package spout

import (
	"log"
)

type  Spout struct {
	Id int
	Data string
	
}


func (*Spout) NextTuple(){
	log.Println("all is okay")
}
