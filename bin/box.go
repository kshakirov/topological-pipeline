package main

func NewBox(id string, fn BoxFunc) *Box {
	return &Box{
		ID:          id,
		UserFunc:    fn,

	}
}

// Start запускает схемотехнику Ящика
func (b *Box) Start(msg Set ) Set{
	
	res := b.UserFunc(msg)
	return res
	

}

