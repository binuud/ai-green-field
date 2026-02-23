package btensor

type BTensor struct {
	Data []float32
	cols int
	rows int
}

func NewFromArange(start, stop, step float32) *BTensor {
	data := Arange(start, stop, step)
	bTensor := &BTensor{
		Data: data,
		cols: len(data),
		rows: 1,
	}
	return bTensor
}
