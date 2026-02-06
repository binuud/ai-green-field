package neuralNetwork

import (
	"fmt"
	"image/color"
	"math"
)

const (
	lr     = 0.01
	epochs = 3000
)

type neuralNetwork struct {
	config neuralNetworkConfig
}

type neuralNetworkConfig struct {
	inNeurons     int
	outNeurons    int
	hiddenNeurons int
	numEpochs     int
	learningRate  float64
}

type LinearParams struct {
	Weight float64
	Bias   float64
}

func NewNeuralNetwork(config neuralNetworkConfig) *neuralNetwork {
	return &neuralNetwork{config: config}
}

func UpdateLinearParams(p *LinearParams, gradientM float64, gradientC float64, lr float64) {
	p.Weight -= lr * gradientM
	p.Bias -= lr * gradientC
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

func Fit(xTrain []float64, yTrain []float64, xTest []float64, yTest []float64, p *LinearParams) {

	graph := NewPlotter()
	graph.addPoints(xTrain, yTrain, color.RGBA{0, 0, 255, 255})
	graph.addPoints(xTest, yTest, color.RGBA{0, 255, 0, 255})
	graph.saveGraph()

	for i := 0; i < epochs; i++ {

		predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		loss := CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		UpdateLinearParams(p, gradM, gradC, lr)

		if i%50 == 0 {
			predictedTest := ApplyLinearEquation(xTest, p.Weight, p.Bias)
			graph.addPoints(xTest, predictedTest, color.RGBA{255, 0, 0, 255})
			fmt.Printf("\n Loss: %f Weight: %f Bias:%f", loss, p.Weight, p.Bias)
		}

	}

	graph.saveGraph()

}
