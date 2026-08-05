package main

//оператор внешнего выбора -External Choice
type SplitDispatcher interface{
	WriteWithChoice(in Set) 
}


//оператор чередования Interleaving
type SplitBuffer interface {
	Interleave(out Set)
}

//параллельная композиция
type SplitNode interface{
	WriteWithChoice(in Set)
	Interleave(out Set)
	
}


type LocalSplitDispatcher struct{
	currentIndex int;
	channelsNumber int;
	InChan chan Set
	BoxesChans chan[] Set
	//Wire который надо передать сорсу
	
}

func (lsd *LocalSplitDispatcher) WriteWithChoice (in Set, boxesChans chan[] Set){
	//take current index пиши по модулю 
}

type LocalSplitBuffer struct{
	buffer []Set
	SinkChan Wire
	
}

func (lsp *LocalSplitBuffer) Interleave(out Set){
	// push all to the buffer the buffer is somehow connected to Wire
}

type LocalSplitNode struct{
	Dispatcher LocalSplitDispatcher
	Nodes []ComputeNode
	Buffer LocalSplitBuffer
}

//func (lsn *LocalSplitNode) Writer (Set) 
