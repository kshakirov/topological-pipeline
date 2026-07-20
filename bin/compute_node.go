package main
import("log")

func (s *ComputeNode) AddInOut(i *Box, o *Box){
	s.OutBox = o
	s.InBox = i
}

func (s *ComputeNode) Prep(){
	s.InChan = make(chan Set, 10)
	s.OutChan = make(chan Set, 10)
}


func (s *ComputeNode) WireIn(){
	go func() {
		for msg := range s.InChan {
			log.Printf("ComputeNode[%d]: The msg arrived\n",s.Id );
			res := s.InBox.UserFunc(msg)
			log.Printf("InBox [%s] applied to the msg, result is [%v]\n", s.InBox.ID, res)
			res = s.OutBox.UserFunc(res)
			log.Printf("OutBox [%s] applied to the msg result is [%v]\n", s.OutBox.ID, res)
			s.OutChan <- res
		}
	}()

}
