package main

type Status byte

const (
	StatusUnknown Status = iota // 0 — зарезервировано для неизвестного/неустановленного
	Start                    // 1
	Magic                 // 2
	Length
	Frame
)



func Framer(payload []byte) ( msgs [][]byte,tail []byte){
	//heads:= make([][]byte,0)
	//	current_head := make([]byte,0)
	current_segment_length :=0
	current_length :=0
	frame_length :=0
	status := Magic
	
	for _,p := range(payload) {
		comb := [3]byte{byte(status),p,byte(current_length)}
		switch(comb){
			case [3]byte{byte(Magic), 0x53, 0x0}: {
					current_segment_length += 1
					current_length +=1

			}
			case [3]byte{byte(Magic), 0x53, 0x1}: {
					status = Length
					current_length +=1


			}
			

			case [3]byte{byte(Length), p, 0x2}: {
				current_length +=1
				frame_length = int(p << 8)

			}
			case [3]byte{byte(Length), p, 0x3}: {
				current_length +=1
				status = Frame
				frame_length = frame_length + int(p)
			}

			default: {
				switch(status){
					case Frame: {
						
					}
					
				}
			}

			
		}
	}
	return nil, nil
	
}
