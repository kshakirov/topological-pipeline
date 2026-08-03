package main

type DemuxNode struct {
	OutBoxes []Box
	InBox *Box
}


func (d *DemuxNode) AddIn(box *Box){
	d.InBox = box
}

func (d *DemuxNode) AddOut(box Box){
	d.OutBoxes = append(d.OutBoxes, box)
}

