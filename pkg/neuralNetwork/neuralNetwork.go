package neuralNetwork

import (
	"fmt"
	"image/color"

	"github.com/sirupsen/logrus"
)

const (
	lr     = 0.01
	epochs = 3000
)

type neuralNetwork struct {
	config *NeuralNetworkConfig
}

type NeuralNetworkConfig struct {
	InNeurons        int
	OutNeurons       int
	HiddenNeurons    int
	NumEpochs        int
	Name             string // Name of the model
	ModelFile        string // Name of the model file for loading and saving model
	LearningRate     float64
	Seed             int64
	LinearRegression *linearRegressionModel
}

type LinearParams struct {
	Weight float64
	Bias   float64
}

func NewNeuralNetwork(config *NeuralNetworkConfig) *neuralNetwork {

	linearRegression := NewLinearRegressionModel(config.Seed)
	config.LinearRegression = linearRegression
	logrus.Infof("\n Model details %v", config)
	return &neuralNetwork{config: config}

}

func NewNeuralNetworkFromModel(ModelFile string) (error, *neuralNetwork) {

	err, config := loadFromModelFile(ModelFile)
	if err != nil {
		return err, nil
	}

	logrus.Infof("\n Model details %v", config)

	return nil, &neuralNetwork{config: config}

}

func (nn *neuralNetwork) GetConfig() *NeuralNetworkConfig {
	return nn.config
}

func (nn *neuralNetwork) LogConfig() {
	logrus.Infof("Model Config %#v", nn.config)
	logrus.Infof("Linear Model%#v", nn.config.LinearRegression.Params)
}

func (nn *neuralNetwork) UpdateLinearParams(gradientM float64, gradientC float64, lr float64) {
	nn.config.LinearRegression.Params.Weight -= lr * gradientM
	nn.config.LinearRegression.Params.Bias -= lr * gradientC
}

func (nn *neuralNetwork) Train(xTrain []float64, yTrain []float64, xTest []float64, yTest []float64) {

	graph := NewPlotter("linearPlot.png")
	graph.addPoints(xTrain, yTrain, color.RGBA{0, 0, 255, 255})
	graph.addPoints(xTest, yTest, color.RGBA{0, 255, 0, 255})
	graph.saveGraph()

	p := nn.config.LinearRegression.Params

	for i := 0; i < epochs; i++ {

		predicted := nn.Predict(xTrain)
		// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		loss := CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		nn.UpdateLinearParams(gradM, gradC, lr)

		if i%50 == 0 {
			predictedTest := ApplyLinearEquation(xTest, p.Weight, p.Bias)
			graph.addPoints(xTest, predictedTest, color.RGBA{255, 0, 0, 255})
			fmt.Printf("\n Loss: %f Weight: %f Bias:%f", loss, p.Weight, p.Bias)
		}

	}

	graph.saveGraph()

}

func (nn *neuralNetwork) Predict(testData []float64) []float64 {
	return nn.config.LinearRegression.Forward(testData)
}
