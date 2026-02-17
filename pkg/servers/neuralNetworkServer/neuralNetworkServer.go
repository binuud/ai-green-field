package neuralnetworkserver

import (
	"context"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcNeuralnetworkserver struct {
	trainingState *protoV1.TrainingState
	protoV1.UnimplementedNeuralNetworkServer
}

func NewGrpcNeuralNetworkServer() *grpcNeuralnetworkserver {

	return &grpcNeuralnetworkserver{
		trainingState: nil,
	}

}

func (s *grpcNeuralnetworkserver) Create(ctx context.Context, in *protoV1.CreateNeuralNetworkRequest) (*protoV1.CreateNeuralNetworkResponse, error) {

	state := in.State
	state.Uuid = uuid.New().String() // #TODO change to uuidv4

	state.CreatedAt = timestamppb.New(time.Now())
	state.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Created Neural Network (%v)", state)

	return &protoV1.CreateNeuralNetworkResponse{
		State: state,
	}, nil
}

func (s *grpcNeuralnetworkserver) Ping(ctx context.Context, in *protoV1.PingNeuralNetworkRequest) (*protoV1.PingNeuralNetworkResponse, error) {

	logrus.Infof("Ping Neural Network")

	return &protoV1.PingNeuralNetworkResponse{
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
