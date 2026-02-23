package btensor

import "math"

func MakeSlice[T any](values ...T) []T {
	return values
}

func Arange(start, stop, step float32) []float32 {
	N := int(math.Ceil(float64((stop - start) / step)))
	result := make([]float32, N)
	for x := range result {
		result[x] = start + step*float32(x)
	}
	return result
}
