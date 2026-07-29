package main

import (
_"log"


)

func (s *ComputeNode) AddInOut(i *Box, o *Box){
	s.OutBox = o
	s.InBox = i
}



func (s *ComputeNode) Prep() {
    // Ящик просто втыкает свои обработчики в провод. 
    // Вся асинхронная магия и горутины запускаются внутри провода.
    s.Wire.WireIn(*s.InBox, *s.OutBox)
}
