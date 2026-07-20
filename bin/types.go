package main

type Set interface{

}

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


type Wire interface{
	Put(msg Set)
	Get() (Set,bool)
	Close()	
}
