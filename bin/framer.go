package main

import (
	"fmt"

)

type Status byte

const (
	StatusUnknown Status = iota // 0 — зарезервировано для неизвестного/неустановленного
	Start                       // 1
	Magic                       // 2
	Length
	Frame
	Tail
)

func Framer(payload []byte) (msgs [][]byte, tail []byte) {


	current_length := 0
	frame_length := 0
	frame_start := 0
	frames := make([][]byte, 0)
	frame_counter := 0
	aux_length :=0
	msg_start := 0
	status := Magic
	

	for idx, p := range payload {
		comb := [3]byte{byte(status), p, byte(current_length)}
		fmt.Printf("idx[%d] = %v\n", idx, comb);
		switch comb {
			
		case [3]byte{byte(Magic), 0x53, 0x0}:
			{
				current_length += 1
				aux_length +=1
				msg_start =idx

			}
		case [3]byte{byte(Magic), 0x54, 0x1}:
			{
				current_length += 1
				status = Length
				aux_length +=1
			

			}

		case [3]byte{byte(Length), p, 0x2}:
			{
				current_length += 1
				frame_length = int(p << 8)
				aux_length += 1
			}
		case [3]byte{byte(Length), p, 0x3}:
			{
				current_length += 1
				status = Frame
				frame_length = frame_length + int(p)
	
				frame_start = idx + 1
				frame_counter = frame_length - 1
				aux_length += 1
			}

		default:
			{
				switch status {
				case Frame:
					{
						current_length += 1
						if frame_counter == 0 {
							status = Magic
							frames = append(frames, payload[frame_start:frame_start + frame_length])
							frame_length = 0
							current_length = 0
							aux_length =0
						}
						frame_counter -= 1
					}

				}
			}

		}
	}
	fmt.Printf("Current frame length is %d\n",frame_length)
	fmt.Printf("Current  length is %d aux length %d \n",current_length,aux_length)

	if len(frames) == 0 {
		frames = nil 
		fmt.Printf("there is no ready frame, current length is %d\n",current_length)
		tail = payload
	}else{
		if current_length - aux_length < frame_length {
			fmt.Printf("there is a  tail its  current length is %d\n",current_length)
			tail =payload[msg_start:]

		}
	}
	fmt.Printf("frames %v\n", frames)
	fmt.Printf("tail %v\n", tail)

	return frames, tail

}
