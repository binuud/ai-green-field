package main

import (
	"fmt"

	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.Info("Init")

	// Create data
	// start := float32(0.0)
	// end := float32(1.0)
	// step := float32(0.02)

	x := bTensor.NewFromArange(0.0, 1.0, 0.02)
	fmt.Println("Result X:", x.Data[:10])

	actualLinearParams := &nn.LinearParams{
		Weight: 0.7,
		Bias:   0.3,
	}

	y := nn.ApplyLinearEquation(x.Data, actualLinearParams.Weight, actualLinearParams.Bias)
	fmt.Println("Result Y:", y[:10])

	// Create train/test split
	train_split := int(0.8 * float32(len(x.Data))) // 80% of data used for training set, 20% for testing
	xTrain, yTrain := x.Data[:train_split], y[:train_split]
	xTest, yTest := x.Data[train_split:], y[train_split:]

	fmt.Printf("\n Training Data len %d, %d", len(xTrain), len(yTrain))
	fmt.Printf("\n Test Data len %d, %d", len(xTest), len(yTest))

	// create random training weights
	model := nn.NewLinearRegressionModel(42.0)

	fmt.Printf("\n Training weights on model: %v", model)

	nn.Fit(xTrain, yTrain, xTest, yTest, model.Params)

}
