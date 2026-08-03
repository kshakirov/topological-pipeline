package main
import("log")


func (s *SourceNode) AddOut(o *Box){
	s.OutBox = o

}




func (s *SourceNode) Start(input Set){
	log.Printf("StartNode: Input = %v\n", s)
	res:= s.OutBox.UserFunc(input)
	log.Printf("StartNode: Output= %v \n", res)
	log.Printf("s.OutChan %v\n", s.Wire)
	lres := Tuple{Fst:res,Snd:false}
	s.Wire.Write(lres.Fst)
	//*(s.InChan) <- res


}


func (s *SourceNode) AddChannel(c Wire){
	s.Wire = c

}
