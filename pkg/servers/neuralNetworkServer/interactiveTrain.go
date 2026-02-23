package neuralnetworkserver

import (
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *grpcNeuralnetworkserver) updateTrainingModel(in *protoV1.Model) {

	if s.Model == nil {
		s.Model = &protoV1.Model{
			Config: nil,
			State: &protoV1.ModelState{
				CurrentEpoch: 0,
				CreatedAt:    timestamppb.New(time.Now()),
				UpdatedAt:    timestamppb.New(time.Now()),
				TrainingLoss: 1,
				TestLoss:     1,
			},
			Uuid: uuid.New().String(),
		}

	}

	if s.Model.Uuid == "" {
		s.Model.Uuid = uuid.New().String()
	}

	config := &protoV1.ModelConfig{
		Name:               in.Config.Name,
		ActivationFunction: in.Config.ActivationFunction,
		EpochBatch:         in.Config.EpochBatch,
		Epochs:             in.Config.Epochs,
		LearningRate:       in.Config.LearningRate,
		NumInputs:          in.Config.NumInputs,
		NumLayers:          in.Config.NumLayers,
		NumOutputs:         in.Config.NumOutputs,
		Regularization:     in.Config.Regularization,
		RegularizationRate: in.Config.RegularizationRate,
	}

	s.Model.Config = config
	s.Model.State.Status = protoV1.ModelState_Pause

}

func (s *grpcNeuralnetworkserver) sendStreamResponse(stream protoV1.NeuralNetwork_InteractiveTrainServer) {
	resp := &protoV1.InteractiveTrainNeuralNetworkResponse{
		Model: s.Model,
	}
	logrus.Infof("\n Sending model ping to client %v", resp.Model.State)
	if err := stream.Send(resp); err != nil {
		logrus.Printf("send error: %v", err)
	}
}

// InteractiveTrain method, for pausing and starting training and peeking and loss data during training
// The model can be retuned too, for inspecting model changes via UI
func (s *grpcNeuralnetworkserver) InteractiveTrain(stream protoV1.NeuralNetwork_InteractiveTrainServer) error {

	var (
		currentNum                   = 1
		mu                           sync.Mutex
		x                            *bTensor.BTensor
		y                            []float32
		xTrain, yTrain, xTest, yTest []float32
		nnModel                      *nn.NeuralNetwork
	)

	// Goroutine to handle continuous streaming of training state, when its running
	// go func() {
	// 	ticker := time.NewTicker(1000 * time.Millisecond)
	// 	defer ticker.Stop()

	// 	for {
	// 		mu.Lock()
	// 		// send updates every 1000 msec, only if the training is in running state
	// 		if model != nil && model.State != nil {
	// 			resp := &protoV1.InteractiveTrainNeuralNetworkResponse{
	// 				Model: model,
	// 			}
	// 			logrus.Infof("\n Sending model ping to client %v", resp.Model.State)
	// 			if err := stream.Send(resp); err != nil {
	// 				logrus.Printf("send error: %v", err)
	// 				return
	// 			}
	// 		}

	// 		mu.Unlock()

	// 		<-ticker.C
	// 	}
	// }()

	// Handle client control messages
	for {
		in, err := stream.Recv()
		if err != nil {
			logrus.Printf("recv error: %v", err)
			return err
		} else {

			mu.Lock()
			switch in.Action.Action {

			case protoV1.TrainingAction_New:

				logrus.Printf("\n Initializing training with fresh set (%v)", in.Model)
				s.updateTrainingModel(in.Model)

				x = bTensor.NewFromArange(0.0, 1.0, 0.02)
				logrus.Println("Result X:", x.Data[:10])

				actualLinearParams := &protoV1.LinearRegressionModel{
					Weight: 1.05,
					Bias:   0.95,
				}

				y = nn.ApplyLinearEquation(x.Data, actualLinearParams.Weight, actualLinearParams.Bias)
				logrus.Println("Result Y:", y[:10])

				// Create train/test split
				train_split := int(0.8 * float32(len(x.Data))) // 80% of data used for training set, 20% for testing
				xTrain, yTrain = x.Data[:train_split], y[:train_split]
				xTest, yTest = x.Data[train_split:], y[train_split:]

				logrus.Printf("\n Training Data len %d, %d", len(xTrain), len(yTrain))
				logrus.Printf("\n Test Data len %d, %d", len(xTest), len(yTest))

				// create random training weights
				nnModel = nn.NewNeuralNetworkFromModel(s.Model)
				nnModel.LogConfig()
				s.sendStreamResponse(stream)

			case protoV1.TrainingAction_Start:
				// Training resumed after pause,
				// Train for another epoch batch, and return the updated model with new test and training loss
				logrus.Println("Starting training")
				s.Model.State.Status = protoV1.ModelState_Running
				nnModel.InteractiveTrain(xTrain, yTrain, xTest, yTest)
				nnModel.LogConfig()
				s.Model.State.CurrentEpoch += s.Model.Config.EpochBatch

				s.Model.State = s.Model.State
				if s.Model.State.CurrentEpoch >= s.Model.Config.EpochBatch {
					s.Model.State.Status = protoV1.ModelState_Completed
				} else {
					s.Model.State.Status = protoV1.ModelState_Pause
				}
				s.sendStreamResponse(stream)

			case protoV1.TrainingAction_Pause:
				logrus.Println("Pausing training")
				s.Model.State.Status = protoV1.ModelState_Pause
				logrus.Printf("Stream paused at number %d", currentNum)

			case protoV1.TrainingAction_Stop:
				s.Model.State.Status = protoV1.ModelState_Stopped
				logrus.Printf("Stream paused at number %d", currentNum)

			default:
				logrus.Printf("Unknown command: %s", in.Action.Action)

			}
			mu.Unlock()
		}

	}
}
