package btensor

import "math"

func MakeSlice[T any](values ...T) []T {
	return values
}

func linearEquation(weight float64, x float64, bias float64) float64 {
	return (weight*x + bias)
}

func Arange(start, stop, step float64) []float64 {
	N := int(math.Ceil((stop - start) / step))
	result := make([]float64, N)
	for x := range result {
		result[x] = start + step*float64(x)
	}
	return result
}
