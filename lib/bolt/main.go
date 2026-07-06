package bolt

import (

	"log"
	tuple "gostorm.org/go_storm/lib/tuple"
)

type Bolt struct{

	Id int

}


func (b *Bolt) Execute(t *tuple.Tuple){
	log.Println("Bolt is okay")
	log.Printf("from tuple %s \n",t.Get(0) )
}
