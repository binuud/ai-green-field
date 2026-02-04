package neuralNetwork

import "math/rand"

type linearRegressionModel struct {
	Params *LinearParams
}

func NewLinearRegressionModel(seed int64) *linearRegressionModel {

	r := rand.New(rand.NewSource(seed))
	p := &LinearParams{Weight: r.Float64(), Bias: r.Float64()}
	return &linearRegressionModel{
		Params: p,
	}
}

func (l *linearRegressionModel) Forward(xData []float64) []float64 {
	return ApplyLinearEquation(xData, l.Params.Weight, l.Params.Bias)
}
