package main

//оператор внешнего выбора -External Choice
type SplitDispatcher interface{
	Write(in Set) 
}


//оператор чередования Interleaving
type SplitBuffer interface {
	Queue(out Set)
}

//параллельная композиция
type SplitNode interface{
	Write(in Set)
	Queue(out Set)
	
}
