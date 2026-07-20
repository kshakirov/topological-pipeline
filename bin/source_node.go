package main
import("log")

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


func (s *SourceNode) AddChannel(c *chan Set){
	s.InChan = c

}



func wireStartNode(s *SourceNode,c *ComputeNode){
	s.AddChannel(&c.InChan)
	
}
