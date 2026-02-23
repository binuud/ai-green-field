package neuralNetwork

import "math"

func linearEquation(x float32, weight float32, bias float32) float32 {
	return (weight*x + bias)
}

func ApplyLinearEquation(xData []float32, weight float32, bias float32) []float32 {
	result := make([]float32, len(xData))
	for i, x := range xData {
		result[i] = linearEquation(x, weight, bias)
	}
	return result
}

// MSE or Mean Squared Error
func CalculateLoss(actuals []float32, predictions []float32) float32 {
	var diff float64
	for i := range actuals {
		d := predictions[i] - actuals[i]
		diff += math.Pow(float64(d), 2.0)
	}
	avergageLoss := diff / float64(len(actuals))
	return float32(avergageLoss)
}

func CalcGradientWeight(orig_x []float32, orig_y []float32, predicted_y []float32) float32 {
	var diff float32
	for i, x := range orig_x {
		diff += (predicted_y[i] - orig_y[i]) * x
	}
	return float32(2 * (diff / float32(len(orig_x))))
}

func CalcGradientBias(orig []float32, predicted []float32) float32 {

	var diff float32
	for i := range orig {
		diff += (predicted[i] - orig[i])
	}
	return float32(2 * (diff / float32(len(orig))))

}

func CalculateGradients(xTrain []float32, yTrain []float32, predicted []float32) (gradM float32, gradC float32) {

	gradM = CalcGradientWeight(xTrain, yTrain, predicted)
	gradC = CalcGradientBias(yTrain, predicted)
	return

}
