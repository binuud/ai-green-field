package neuralnetworkserver

import (
	"context"
	"fmt"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	datalayer "github.com/binuud/ai-green-field/pkg/dataLayer"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
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

	s.updateTrainingModel(in.Model)

	logrus.Infof("Created Neural Network (%v)", s.Model)

	nNet := nn.NewNeuralNetworkFromModel(s.Model)
	err := nNet.Save("")
	if err != nil {
		logrus.Errorf("Cannot save model %v", err)
	}

	return &protoV1.CreateNeuralNetworkResponse{
		Model: s.Model,
	}, nil
}

func (s *grpcNeuralnetworkserver) Save(ctx context.Context, in *protoV1.SaveNeuralNetworkRequest) (*protoV1.SaveNeuralNetworkResponse, error) {

	model := in.Model

	model.State.UpdatedAt = timestamppb.New(time.Now())

	nNet := nn.NewNeuralNetworkFromModel(s.Model)
	err := nNet.Save("")
	if err != nil {
		logrus.Errorf("Cannot save model %v", err)
	} else {
		logrus.Infof("Saved Neural Network (%v)", model)
	}

	return &protoV1.SaveNeuralNetworkResponse{
		Model: model,
	}, err
}

func (s *grpcNeuralnetworkserver) Load(ctx context.Context, in *protoV1.LoadNeuralNetworkRequest) (*protoV1.LoadNeuralNetworkResponse, error) {

	model := in.Model

	model.State.UpdatedAt = timestamppb.New(time.Now())

	nNet, err := nn.NewNeuralNetworkFromModelFile(s.Model.Uuid)
	if err != nil {
		logrus.Errorf("Cannot load model %s", in.Model.Uuid)
		return nil, err
	}

	logrus.Infof("Loaded Neural Network (%v)", nNet.Model)

	return &protoV1.LoadNeuralNetworkResponse{
		Model: nNet.Model,
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
	nNet := nn.NewNeuralNetwork(model.Config)

	nNet.LogConfig()
	nNet.Train(xTrain, yTrain, xTest, yTest)
	nNet.LogConfig()
	//check model loss with test data
	predicted := nNet.Predict(xTest)
	// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
	loss := nn.CalculateLoss(yTest, predicted)
	logrus.Printf("\n Loss on test data %f", loss)

	nNet.Model.State.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Training Neural Network (%v)", nNet.Model)

	return &protoV1.TrainNeuralNetworkResponse{
		Model: nNet.Model,
	}, nil
}

func (s *grpcNeuralnetworkserver) List(ctx context.Context, in *protoV1.ListNeuralNetworkRequest) (*protoV1.ListNeuralNetworkResponse, error) {

	logrus.Infof("List Neural Network")
	fileList, err := datalayer.ListFilesNonRecursive(nn.GetNNStoreDirectory())
	if err != nil {
		logrus.Errorf("Cannot list model %s", in.Model.Uuid)
		return nil, err
	}
	logrus.Infof("File list %v", fileList)

	return &protoV1.ListNeuralNetworkResponse{
		Files: fileList,
	}, nil
}

func (s *grpcNeuralnetworkserver) Ping(ctx context.Context, in *protoV1.PingNeuralNetworkRequest) (*protoV1.PingNeuralNetworkResponse, error) {

	logrus.Infof("Ping Neural Network")

	return &protoV1.PingNeuralNetworkResponse{
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
