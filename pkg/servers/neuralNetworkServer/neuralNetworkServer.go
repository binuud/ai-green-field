package neuralnetworkserver

import (
	"context"
	"fmt"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcNeuralnetworkserver struct {
	Model *protoV1.Model
	protoV1.UnimplementedNeuralNetworkServer
}

func NewGrpcNeuralNetworkServer() *grpcNeuralnetworkserver {

	return &grpcNeuralnetworkserver{
		Model: nil,
	}

}

func (s *grpcNeuralnetworkserver) Create(ctx context.Context, in *protoV1.CreateNeuralNetworkRequest) (*protoV1.CreateNeuralNetworkResponse, error) {

	model := in.Model
	model.Uuid = uuid.New().String() // #TODO change to uuidv4

	model.State.CreatedAt = timestamppb.New(time.Now())
	model.State.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Created Neural Network (%v)", model)

	return &protoV1.CreateNeuralNetworkResponse{
		Model: model,
	}, nil
}

func (s *grpcNeuralnetworkserver) Save(ctx context.Context, in *protoV1.SaveNeuralNetworkRequest) (*protoV1.SaveNeuralNetworkResponse, error) {

	model := in.Model
	model.Uuid = uuid.New().String() // #TODO change to uuidv4

	model.State.CreatedAt = timestamppb.New(time.Now())
	model.State.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Saved Neural Network (%v)", model)

	return &protoV1.SaveNeuralNetworkResponse{
		Model: model,
	}, nil
}

func (s *grpcNeuralnetworkserver) Train(ctx context.Context, in *protoV1.TrainNeuralNetworkRequest) (*protoV1.TrainNeuralNetworkResponse, error) {

	model := in.Model

	x := bTensor.NewFromArange(0.0, 1.0, 0.02)
	fmt.Println("Result X:", x.Data[:10])

	actualLinearParams := &protoV1.LinearRegressionModel{
		Weight: 1.05,
		Bias:   0.95,
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
	nnModel := nn.NewNeuralNetwork(model.Config)

	nnModel.LogConfig()
	nnModel.Train(xTrain, yTrain, xTest, yTest)
	nnModel.LogConfig()
	//check model loss with test data
	predicted := nnModel.Predict(xTest)
	// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
	loss := nn.CalculateLoss(yTest, predicted)
	logrus.Printf("\n Loss on test data %f", loss)

	nnModel.Model.State.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Training Neural Network (%v)", nnModel.Model)

	return &protoV1.TrainNeuralNetworkResponse{
		Model: nnModel.Model,
	}, nil
}

func (s *grpcNeuralnetworkserver) Ping(ctx context.Context, in *protoV1.PingNeuralNetworkRequest) (*protoV1.PingNeuralNetworkResponse, error) {

	logrus.Infof("Ping Neural Network")

	return &protoV1.PingNeuralNetworkResponse{
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
