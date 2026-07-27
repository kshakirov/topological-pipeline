package main
import("log")


func (s *SourceNode) AddOut(o *Box){
	s.OutBox = o

}




func (s *SourceNode) Start(input Set){
	log.Printf("StartNode: Input = %v\n", s)
	res:= s.OutBox.UserFunc(input)
	log.Printf("StartNode: Output= %v \n", res)
	lres := Tuple{Fst:res,Snd:false}
	(*s.OutChan).Write(lres, true)
	//*(s.InChan) <- res


}


func (s *SourceNode) AddChannel(c *OutWire){
	s.OutChan = c

}



func wireStartNode(s *SourceNode,c *ComputeNode){
	s.AddChannel(&c.OutChan)
	
}
