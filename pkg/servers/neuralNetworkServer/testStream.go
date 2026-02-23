package neuralnetworkserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestStream method for testing integration with UI and other bidirectional client
// This sends a counter in numOutputs, handles Pause and Continue action messages
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
				Model: &protoV1.Model{
					State: &protoV1.ModelState{
						CurrentEpoch: int32(currentNum),
						UpdatedAt:    timestamppb.New(time.Now()),
						Status:       protoV1.ModelState_Pause,
					},
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
			logrus.Printf("recv error: %v", err)
			return err
		} else {
			mu.Lock()
			switch in.Action.Action {
			case protoV1.TrainingAction_Start:
				paused = false
				logrus.Printf("Stream resumed at number %d", currentNum)
			case protoV1.TrainingAction_Pause:
				paused = true
				logrus.Printf("Stream paused at number %d", currentNum)
			default:
				logrus.Printf("Unknown command: %s", in.Action.Action)
			}
			mu.Unlock()
		}

	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

// TestStreamWebsocketHandler listens on a websocket and handles the bidirectional stream for the grpc server
// this is for quickly testing any client, that supports bidirectional messages
func TestStreamWebsocketHandler(w http.ResponseWriter, r *http.Request) {
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
	grpcStream, err := client.TestStream(ctx)
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
			var inMsg protoV1.TestStreamNeuralNetworkRequest
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
