package neuralNetwork

import (
	"image/color"

	"github.com/sirupsen/logrus"
)

const (
	lr     = 0.01
	epochs = 3000
)

type NeuralNetwork struct {
	config *NeuralNetworkConfig
}

type NeuralNetworkConfig struct {
	InNeurons        int
	OutNeurons       int
	HiddenNeurons    int
	NumEpochs        int
	EpochBatch       int    // number of epochs to run and return
	Name             string // Name of the model
	ModelFile        string // Name of the model file for loading and saving model
	LearningRate     float64
	Seed             int64
	TestLoss         float64
	TrainingLoss     float64
	LinearRegression *linearRegressionModel
}

type LinearParams struct {
	Weight float64
	Bias   float64
}

func NewNeuralNetwork(config *NeuralNetworkConfig) *NeuralNetwork {

	linearRegression := NewLinearRegressionModel(config.Seed)
	config.LinearRegression = linearRegression
	logrus.Infof("\n Model details %v", config)
	return &NeuralNetwork{config: config}

}

func NewNeuralNetworkFromModel(ModelFile string) (error, *NeuralNetwork) {

	err, config := loadFromModelFile(ModelFile)
	if err != nil {
		return err, nil
	}

	logrus.Infof("\n Model details %v", config)

	return nil, &NeuralNetwork{config: config}

}

func (nn *NeuralNetwork) GetConfig() *NeuralNetworkConfig {
	return nn.config
}

func (nn *NeuralNetwork) LogConfig() {
	logrus.Infof("\n Model Config %#v", nn.config)
	logrus.Infof("\n Linear Model%#v", nn.config.LinearRegression.Params)
}

func (nn *NeuralNetwork) UpdateLinearParams(gradientM float64, gradientC float64, lr float64) {
	nn.config.LinearRegression.Params.Weight -= lr * gradientM
	nn.config.LinearRegression.Params.Bias -= lr * gradientC
}

func (nn *NeuralNetwork) Train(xTrain []float64, yTrain []float64, xTest []float64, yTest []float64) {

	graph := NewPlotter("linearPlot.png")
	graph.addPoints(xTrain, yTrain, color.RGBA{0, 0, 255, 255})
	graph.addPoints(xTest, yTest, color.RGBA{0, 255, 0, 255})
	graph.saveGraph()

	p := nn.config.LinearRegression.Params

	for i := range epochs {

		predicted := nn.Predict(xTrain)
		// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		nn.config.TrainingLoss = CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		nn.UpdateLinearParams(gradM, gradC, lr)

		if i%nn.config.EpochBatch == 0 {
			predictedTest := ApplyLinearEquation(xTest, p.Weight, p.Bias)
			graph.addPoints(xTest, predictedTest, color.RGBA{255, 0, 0, 255})
			logrus.Printf("\n Loss: %f Weight: %f Bias:%f", nn.config.TrainingLoss, p.Weight, p.Bias)
		}

	}

	graph.saveGraph()

}

func (nn *NeuralNetwork) InteractiveTrain(xTrain []float64, yTrain []float64, xTest []float64, yTest []float64, epochBatch int) {

	p := nn.config.LinearRegression.Params

	for range epochBatch {

		predicted := nn.Predict(xTrain)
		// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		nn.config.TrainingLoss = CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		nn.UpdateLinearParams(gradM, gradC, lr)

	}
	logrus.Printf("\n Loss: %f Weight: %f Bias:%f", nn.config.TrainingLoss, p.Weight, p.Bias)

}

func (nn *NeuralNetwork) Predict(testData []float64) []float64 {
	return nn.config.LinearRegression.Forward(testData)
}
