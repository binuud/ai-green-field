package btensor

type BTensor struct {
	Data []float64
	cols int
	rows int
}

func NewFromArange(start, stop, step float64) *BTensor {
	data := Arange(start, stop, step)
	bTensor := &BTensor{
		Data: data,
		cols: len(data),
		rows: 1,
	}
	return bTensor
}

func (b *BTensor) ApplyLinearEquation(weight float64, bias float64) []float64 {
	result := make([]float64, b.cols)
	for i, x := range b.Data {
		result[i] = linearEquation(weight, x, bias)
	}
	return result
}
