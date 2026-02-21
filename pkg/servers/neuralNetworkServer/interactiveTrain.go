package neuralnetworkserver

import (
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/sirupsen/logrus"
)

func (s *grpcNeuralnetworkserver) newTrainingState(inState *protoV1.TrainingState) *protoV1.TrainingState {

	state := &protoV1.TrainingState{
		State: protoV1.TrainingState_Pause,
	}

	state.State = protoV1.TrainingState_New
	state.ActivationFunction = inState.ActivationFunction
	state.CurrentEpoch = 0
	state.EpochBatch = inState.EpochBatch
	state.Epochs = inState.Epochs
	state.LearningRate = inState.LearningRate
	state.NumInputs = inState.NumInputs
	state.NumLayers = inState.NumLayers
	state.NumOutputs = inState.NumOutputs
	state.Regularization = inState.Regularization
	state.RegularizationRate = inState.RegularizationRate

	return state

}

// InteractiveTrain method, for pausing and starting training and peeking and loss data during training
// The model can be retuned too, for inspecting model changes via UI
func (s *grpcNeuralnetworkserver) InteractiveTrain(stream protoV1.NeuralNetwork_InteractiveTrainServer) error {

	var (
		currentNum = 1
		mu         sync.Mutex
		state      = &protoV1.TrainingState{
			State: protoV1.TrainingState_Pause,
		}
		x                            *bTensor.BTensor
		y                            []float64
		xTrain, yTrain, xTest, yTest []float64
		model                        *nn.NeuralNetwork
	)

	// Goroutine to handle continuous streaming of training state, when its running
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for {
			mu.Lock()
			// send updates every 1000 msec, only if the training is in running state
			if state.State == protoV1.TrainingState_Running {
				resp := &protoV1.InteractiveTrainNeuralNetworkResponse{
					State: state,
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

				logrus.Printf("\n Initializing training with fresh set (%v)", in.State)
				state = s.newTrainingState(in.State)

				x = bTensor.NewFromArange(0.0, 1.0, 0.02)
				logrus.Println("Result X:", x.Data[:10])

				actualLinearParams := &nn.LinearParams{
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
				model = nn.NewNeuralNetwork(&nn.NeuralNetworkConfig{
					Name:          "LinearRegression",
					LearningRate:  float64(state.LearningRate),
					NumEpochs:     int(state.Epochs),
					EpochBatch:    int(state.EpochBatch),
					InNeurons:     int(state.NumInputs),
					OutNeurons:    int(state.NumOutputs),
					HiddenNeurons: int(state.NumLayers),
					Seed:          42.0,
				})

				model.LogConfig()

			case protoV1.TrainingAction_Start:

				logrus.Println("Starting training")
				state.State = protoV1.TrainingState_Running
				model.LogConfig()
				model.InteractiveTrain(xTrain, yTrain, xTest, yTest, int(state.EpochBatch))
				model.LogConfig()
				state.CurrentEpoch += state.EpochBatch
				if state.CurrentEpoch >= state.Epochs {
					state.State = protoV1.TrainingState_Completed
				} else {
					state.State = protoV1.TrainingState_Pause
				}

			case protoV1.TrainingAction_Pause:
				logrus.Println("Pausing training")
				state.State = protoV1.TrainingState_Pause
				logrus.Printf("Stream paused at number %d", currentNum)

			case protoV1.TrainingAction_Stop:
				state.State = protoV1.TrainingState_Stopped
				logrus.Printf("Stream paused at number %d", currentNum)

			default:
				logrus.Printf("Unknown command: %s", in.Action.Action)

			}
			mu.Unlock()
		}

	}
}
