package main

type PrefixGeneratorFuncTCP func()
type SinkFuncTCP func()




type PrefixGeneratorTCP struct {
	Func PrefixGeneratorFuncTCP

}

func (generator *PrefixGeneratorTCP) Start() {
	go generator.Func()
}


type SinkTCP struct {
	Func   SinkFuncTCP
	InChan chan byte
}


func (sink * SinkTCP) Consume(){

	sink.Func();
}
