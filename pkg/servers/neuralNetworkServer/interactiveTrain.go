package neuralnetworkserver

import (
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/sirupsen/logrus"
)

func (s *grpcNeuralnetworkserver) newTrainingState(inConfig *protoV1.ModelConfig) *protoV1.Model {

	state := &protoV1.ModelState{
		Status:       protoV1.ModelState_Pause,
		CurrentEpoch: 0,
	}

	config := &protoV1.ModelConfig{
		Name:               inConfig.Name,
		ActivationFunction: inConfig.ActivationFunction,
		EpochBatch:         inConfig.EpochBatch,
		Epochs:             inConfig.Epochs,
		LearningRate:       inConfig.LearningRate,
		NumInputs:          inConfig.NumInputs,
		NumLayers:          inConfig.NumLayers,
		NumOutputs:         inConfig.NumOutputs,
		Regularization:     inConfig.Regularization,
		RegularizationRate: inConfig.RegularizationRate,
	}

	return &protoV1.Model{
		State:  state,
		Config: config,
	}

}

// InteractiveTrain method, for pausing and starting training and peeking and loss data during training
// The model can be retuned too, for inspecting model changes via UI
func (s *grpcNeuralnetworkserver) InteractiveTrain(stream protoV1.NeuralNetwork_InteractiveTrainServer) error {

	var (
		currentNum                   = 1
		mu                           sync.Mutex
		model                        = &protoV1.Model{}
		x                            *bTensor.BTensor
		y                            []float32
		xTrain, yTrain, xTest, yTest []float32
		nnModel                      *nn.NeuralNetwork
	)

	// Goroutine to handle continuous streaming of training state, when its running
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for {
			mu.Lock()
			// send updates every 1000 msec, only if the training is in running state
			if model.State.Status == protoV1.ModelState_Running {
				resp := &protoV1.InteractiveTrainNeuralNetworkResponse{
					Model: model,
				}
				if err := stream.Send(resp); err != nil {
					logrus.Printf("send error: %v", err)
					return
				}
			}

			mu.Unlock()

			<-ticker.C
		}
	}()

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
				model = s.newTrainingState(in.Model.Config)

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
				nnModel = nn.NewNeuralNetwork(model.Config)
				nnModel.LogConfig()

			case protoV1.TrainingAction_Start:

				logrus.Println("Starting training")
				model.State.Status = protoV1.ModelState_Running
				nnModel.LogConfig()
				nnModel.InteractiveTrain(xTrain, yTrain, xTest, yTest)
				nnModel.LogConfig()
				model.State.CurrentEpoch += model.Config.EpochBatch
				if model.State.CurrentEpoch >= model.Config.EpochBatch {
					model.State.Status = protoV1.ModelState_Completed
				} else {
					model.State.Status = protoV1.ModelState_Pause
				}

			case protoV1.TrainingAction_Pause:
				logrus.Println("Pausing training")
				model.State.Status = protoV1.ModelState_Pause
				logrus.Printf("Stream paused at number %d", currentNum)

			case protoV1.TrainingAction_Stop:
				model.State.Status = protoV1.ModelState_Stopped
				logrus.Printf("Stream paused at number %d", currentNum)

			default:
				logrus.Printf("Unknown command: %s", in.Action.Action)

			}
			mu.Unlock()
		}

	}
}
