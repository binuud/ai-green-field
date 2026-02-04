package neuralNetwork

func linearEquation(x float64, weight float64, bias float64) float64 {
	return (weight*x + bias)
}

func ApplyLinearEquation(xData []float64, weight float64, bias float64) []float64 {
	result := make([]float64, len(xData))
	for i, x := range xData {
		result[i] = linearEquation(x, weight, bias)
	}
	return result
}
