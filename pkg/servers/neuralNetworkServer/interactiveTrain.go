package neuralnetworkserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	bTensor "github.com/binuud/ai-green-field/pkg/bTensor"
	nn "github.com/binuud/ai-green-field/pkg/neuralNetwork"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

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

				state.State = protoV1.TrainingState_New
				state.ActivationFunction = in.State.ActivationFunction
				state.CurrentEpoch = 0
				state.EpochBatch = in.State.EpochBatch
				state.Epochs = in.State.Epochs
				state.LearningRate = in.State.LearningRate
				state.NumInputs = in.State.NumInputs
				state.NumLayers = in.State.NumLayers
				state.NumOutputs = in.State.NumOutputs
				state.Regularization = in.State.Regularization
				state.RegularizationRate = in.State.RegularizationRate

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

// InteractiveTrainWebsocketHandler listens on a websocket and handles the bidirectional stream for the grpc server
// this is for quickly testing any client, that supports bidirectional messages
func InteractiveTrainWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Printf("\n Failed to upgrade websocket: %v", err)
		return
	}
	defer ws.Close()

	// Connect to gRPC server
	grpcConn, err := grpc.NewClient(
		"localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Wait for connection to be ready
	)

	if err != nil {
		logrus.Errorf("Cannot dial grpc %v", err)
		return
	}
	defer grpcConn.Close()

	client := protoV1.NewNeuralNetworkClient(grpcConn)

	// Create context with cancelation
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Start gRPC stream
	grpcStream, err := client.InteractiveTrain(ctx)
	if err != nil {
		logrus.Errorf("Cannot open TestStream %v", err)
		return
	}

	var wg sync.WaitGroup

	// Goroutine 1: Receive from gRPC stream → Send to WebSocket
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			outMsg, err := grpcStream.Recv()
			if err != nil {
				logrus.Printf("gRPC stream ended: %v", err)
				return
			}
			logrus.Infof("Data grpc -> websocket: %+v\n", outMsg)
			if err := ws.WriteJSON(outMsg); err != nil {
				logrus.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}()

	// Goroutine 2: Receive from WebSocket → Send to gRPC stream
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer grpcStream.CloseSend()

		for {

			// Read message from WebSocket client
			_, msg, err := ws.ReadMessage()
			if err != nil {
				logrus.Errorf("Read from websocket error: %v", err)
				break
			}

			// read StreamRequest from websocket
			var inMsg protoV1.InteractiveTrainNeuralNetworkRequest
			err = protojson.Unmarshal(msg, &inMsg)
			// err = json.Unmarshal(msg, &data)
			if err != nil {
				logrus.Errorf("Error unmarshaling JSON: %v", err)
				continue
			}

			// send via grpc client
			logrus.Infof("Data websocket -> grpc: %+v\n", inMsg)
			if err := grpcStream.Send(&inMsg); err != nil {
				logrus.Printf("gRPC send error: %v", err)
				return
			}
		}
	}()

	wg.Wait()

}
