package neuralNetwork

import (
	"image/color"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	lr     = 0.01
	epochs = 3000
)

type NeuralNetwork struct {
	Model *protoV1.Model
}

func NewNeuralNetwork(config *protoV1.ModelConfig) *NeuralNetwork {

	linearRegression := NewLinearRegressionModel(config.Seed)

	logrus.Infof("\n Model details %v", config)
	return &NeuralNetwork{Model: &protoV1.Model{
		Config: config,
		State: &protoV1.ModelState{
			CurrentEpoch: 0,
			CreatedAt:    timestamppb.New(time.Now()),
			UpdatedAt:    timestamppb.New(time.Now()),
		},
		LinearModel: linearRegression,
	}}

}

func NewNeuralNetworkFromModel(model *protoV1.Model) *NeuralNetwork {

	logrus.Infof("\n Model details %v", model)
	linearRegression := NewLinearRegressionModel(model.Config.Seed)
	model.LinearModel = linearRegression

	return &NeuralNetwork{Model: model}

}

func NewNeuralNetworkFromModelFile(ModelFile string) (*NeuralNetwork, error) {

	err, model := loadFromModelFile(ModelFile)
	if err != nil {
		return nil, err
	}

	logrus.Infof("\n Model details %v", model)

	return &NeuralNetwork{Model: model}, nil

}

func (nn *NeuralNetwork) GetConfig() *protoV1.ModelConfig {
	return nn.Model.Config
}

func (nn *NeuralNetwork) LogConfig() {
	logrus.Infof("\n Model Config (%s) -- %v", nn.Model.Uuid, nn.Model.Config)
	logrus.Infof("\n Model State %v", nn.Model.State)
	logrus.Infof("\n Linear Model %v", nn.Model.LinearModel)
}

func (nn *NeuralNetwork) UpdateLinearParams(gradientM float32, gradientC float32, lr float32) {
	nn.Model.LinearModel.Weight -= lr * gradientM
	nn.Model.LinearModel.Bias -= lr * gradientC
}

func (nn *NeuralNetwork) Train(xTrain []float32, yTrain []float32, xTest []float32, yTest []float32) {

	graph := NewPlotter("linearPlot.png")
	graph.addPoints(xTrain, yTrain, color.RGBA{0, 0, 255, 255})
	graph.addPoints(xTest, yTest, color.RGBA{0, 255, 0, 255})
	graph.saveGraph()

	p := nn.Model.LinearModel

	for i := range epochs {

		predicted := nn.Predict(xTrain)
		// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		nn.Model.State.TrainingLoss = CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		nn.UpdateLinearParams(gradM, gradC, lr)

		if i%int(nn.Model.Config.EpochBatch) == 0 {
			predictedTest := ApplyLinearEquation(xTest, p.Weight, p.Bias)
			graph.addPoints(xTest, predictedTest, color.RGBA{255, 0, 0, 255})
			logrus.Printf("\n Loss: %f Weight: %f Bias:%f", nn.Model.State.TrainingLoss, p.Weight, p.Bias)
		}

	}

	graph.saveGraph()

}

func (nn *NeuralNetwork) InteractiveTrain(xTrain []float32, yTrain []float32, xTest []float32, yTest []float32) {

	p := nn.Model.LinearModel

	for range nn.Model.Config.EpochBatch {

		predicted := nn.Predict(xTrain)
		// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
		nn.Model.State.TrainingLoss = CalculateLoss(yTrain, predicted)

		gradM, gradC := CalculateGradients(xTrain, yTrain, predicted)
		nn.UpdateLinearParams(gradM, gradC, lr)

	}
	nn.CalculateTestLoss(xTest, yTest)
	logrus.Printf("\n Training Loss: %f, Test Loss: %f, Weight: %f Bias:%f", nn.Model.State.TrainingLoss, nn.Model.State.TestLoss, p.Weight, p.Bias)

}

func (nn *NeuralNetwork) CalculateTestLoss(xTest []float32, yTest []float32) {

	predicted := nn.Predict(xTest)
	nn.Model.State.TestLoss = CalculateLoss(yTest, predicted)

}

func (nn *NeuralNetwork) Predict(testData []float32) []float32 {
	return nn.Forward(testData)
}
