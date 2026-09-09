package main

type PrefixGeneratorFuncTCP func()
type SinkFuncTCP func(byte)




type PrefixGeneratorTCP struct {
	Func PrefixGeneratorFuncTCP

}

func (generator *PrefixGeneratorTCP) Start() {
	generator.Func()
}


type SinkTCP struct {
	Func   SinkFuncTCP
	InChan chan byte
}
