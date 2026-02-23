package neuralNetwork

import (
	"math/rand"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
)

type NNLinearRegressionModel struct {
	*protoV1.LinearRegressionModel
}

func NewLinearRegressionModel(seed int64) *protoV1.LinearRegressionModel {

	r := rand.New(rand.NewSource(seed))
	return &protoV1.LinearRegressionModel{Weight: r.Float32(), Bias: r.Float32()}
}

func (l *NeuralNetwork) Forward(xData []float32) []float32 {
	return ApplyLinearEquation(xData, l.Model.LinearModel.Weight, l.Model.LinearModel.Bias)
}
