package neuralnetworkserver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
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

func (s *grpcNeuralnetworkserver) Save(ctx context.Context, in *protoV1.SaveNeuralNetworkRequest) (*protoV1.SaveNeuralNetworkResponse, error) {

	state := in.State
	state.Uuid = uuid.New().String() // #TODO change to uuidv4

	state.CreatedAt = timestamppb.New(time.Now())
	state.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Saved Neural Network (%v)", state)

	return &protoV1.SaveNeuralNetworkResponse{
		State: state,
	}, nil
}

func (s *grpcNeuralnetworkserver) Train(ctx context.Context, in *protoV1.TrainNeuralNetworkRequest) (*protoV1.TrainNeuralNetworkResponse, error) {

	state := in.State
	state.Uuid = uuid.New().String() // #TODO change to uuidv4

	x := bTensor.NewFromArange(0.0, 1.0, 0.02)
	fmt.Println("Result X:", x.Data[:10])

	actualLinearParams := &nn.LinearParams{
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
	model := nn.NewNeuralNetwork(&nn.NeuralNetworkConfig{
		Name:         "LinearRegression",
		LearningRate: float64(state.LearningRate),
		NumEpochs:    int(state.Epochs),
		EpochBatch:   int(state.EpochBatch),
		Seed:         42.0,
	})

	model.LogConfig()
	model.Train(xTrain, yTrain, xTest, yTest)
	model.LogConfig()
	//check model loss with test data
	predicted := model.Predict(xTest)
	// predicted := ApplyLinearEquation(xTrain, p.Weight, p.Bias)
	loss := nn.CalculateLoss(yTest, predicted)
	logrus.Printf("\n Loss on test data %f", loss)

	state.UpdatedAt = timestamppb.New(time.Now())

	logrus.Infof("Training Neural Network (%v)", state)

	return &protoV1.TrainNeuralNetworkResponse{
		State: state,
	}, nil
}

func (s *grpcNeuralnetworkserver) TestStream(stream protoV1.NeuralNetwork_TestStreamServer) error {

	var (
		currentNum = 1
		mu         sync.Mutex
		paused     = false
	)

	// Goroutine to send numbers continuously
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for currentNum <= 100 {
			mu.Lock()
			if paused || currentNum > 100 {
				mu.Unlock()
				continue
			}

			resp := &protoV1.TestStreamNeuralNetworkResponse{
				State: &protoV1.TrainingState{
					NumOutputs: int32(currentNum),
					UpdatedAt:  timestamppb.New(time.Now()),
					State:      protoV1.TrainingState_Pause,
				},
			}
			mu.Unlock()

			if err := stream.Send(resp); err != nil {
				logrus.Printf("send error: %v", err)
				return
			}
			logrus.Infof("Streaming number to client --  %d", currentNum)
			currentNum++

			<-ticker.C
		}
	}()

	// Handle client control messages
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("recv error: %v", err)
			// return err
		} else {
			mu.Lock()
			switch in.Action.Action {
			case protoV1.TrainingAction_Start:
				paused = false
				log.Printf("Stream resumed at number %d", currentNum)
			case protoV1.TrainingAction_Pause:
				paused = true
				log.Printf("Stream paused at number %d", currentNum)
			default:
				log.Printf("Unknown command: %s", in.Action.Action)
			}
			mu.Unlock()
		}

	}
}

func (s *grpcNeuralnetworkserver) Ping(ctx context.Context, in *protoV1.PingNeuralNetworkRequest) (*protoV1.PingNeuralNetworkResponse, error) {

	logrus.Infof("Ping Neural Network")

	return &protoV1.PingNeuralNetworkResponse{
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
