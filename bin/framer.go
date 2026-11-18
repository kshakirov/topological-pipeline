package main

import (
	"slices"
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
	aux_length := 0
	msg_start := 0
	status := Magic

	for idx, p := range payload {
		comb := [3]byte{byte(status), p, byte(current_length)}
		switch comb {

		case [3]byte{byte(Magic), 0x53, 0x0}:
			{
				current_length += 1
				aux_length += 1
				msg_start = idx

			}
		case [3]byte{byte(Magic), 0x54, 0x1}:
			{
				current_length += 1
				status = Length
				aux_length += 1

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
							frames = append(frames, payload[frame_start:frame_start+frame_length])
							frame_length = 0
							current_length = 0
							aux_length = 0
						}
						frame_counter -= 1
					}

				}
			}

		}
	}
	if len(frames) == 0 {
		frames = nil
		tail = payload
	} else {
		if current_length-aux_length < frame_length {
			tail = payload[msg_start:]

		}
	}

	return frames, tail

}

func PrefixFrame(payload []byte) []byte {
	length := len(payload)
	l := [2]byte{byte(length >> 8), byte(length)}
	st := []byte{0x53, 0x54}
	return slices.Concat(st, l[:], payload)
}
