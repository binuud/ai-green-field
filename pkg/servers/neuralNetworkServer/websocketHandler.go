package neuralnetworkserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	protoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

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

	logrus.Infoln("+++ Received new websocket connection")
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
			outStr, err := protojson.Marshal(outMsg)
			if err != nil {
				logrus.Printf("Cannot marshall InteractiveTrainNeuralNetworkResponse to string: %v", err)
				return
			}
			logrus.Infof("Data grpc -> websocket: %+v\n", outMsg)
			// writeJSON emits enum as integers
			// this breaks the parsing at the UI end
			// DO NOT CHANGE, by binu.
			if err := ws.WriteMessage(1, outStr); err != nil {
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
			logrus.Infof("Parsing %s", msg)
			// read StreamRequest from websocket
			var inMsg protoV1.InteractiveTrainNeuralNetworkRequest
			err = protojson.Unmarshal(msg, &inMsg)
			// err = json.Unmarshal(msg, &inMsg)
			if err != nil {
				logrus.Errorf("Error unmarshaling InteractiveTrainNeuralNetworkRequest protojson: %v", err)

				err = json.Unmarshal(msg, &inMsg)
				if err != nil {
					logrus.Errorf("Error unmarshaling InteractiveTrainNeuralNetworkRequest json: %v", err)
				}
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
