package main

func (d *MuxNode) AddIn(box *Box){
	d.OutBox = box
}

func (d *MuxNode) addOut(box Box){
	d.InBoxes = append(d.InBoxes, box)
}

