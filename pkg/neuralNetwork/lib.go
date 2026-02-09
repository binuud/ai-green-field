package neuralNetwork

import "math"

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

// MSE or Mean Squared Error
func CalculateLoss(actuals []float64, predictions []float64) float64 {
	var diff float64
	for i := range actuals {
		d := predictions[i] - actuals[i]
		diff += math.Pow(d, 2.0)
	}
	avergageLoss := diff / float64(len(actuals))
	return avergageLoss
}

func CalcGradientWeight(orig_x []float64, orig_y []float64, predicted_y []float64) float64 {
	var diff float64
	for i, x := range orig_x {
		diff += (predicted_y[i] - orig_y[i]) * x
	}
	return 2 * (diff / float64(len(orig_x)))
}

func CalcGradientBias(orig []float64, predicted []float64) float64 {

	var diff float64
	for i := range orig {
		diff += (predicted[i] - orig[i])
	}
	return 2 * (diff / float64(len(orig)))

}

func CalculateGradients(xTrain []float64, yTrain []float64, predicted []float64) (gradM float64, gradC float64) {

	gradM = CalcGradientWeight(xTrain, yTrain, predicted)
	gradC = CalcGradientBias(yTrain, predicted)
	return

}
