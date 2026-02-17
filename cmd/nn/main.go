package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	nnProtoV1 "github.com/binuud/ai-green-field/gen/go/v1/neuralNetwork"
	neuralnetworkserver "github.com/binuud/ai-green-field/pkg/servers/neuralNetworkServer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	server_grpc_port = flag.Int("grpc_port", 9090, "AI-Green-Field Backend Server GRPC port, no token required unsecured")
	server_http_port = flag.Int("http_port", 9080, "AI-Green-Field Backend Server HTTP port, no token required unsecured")
)

func runHtppServer() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.NewClient(
		fmt.Sprintf("0.0.0.0:%d", *server_grpc_port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwMux := runtime.NewServeMux()

	// Neural Network Service
	err = nnProtoV1.RegisterNeuralNetworkHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	httpMux := http.NewServeMux()
	// attach mqtt bridge to the http mux
	// this is a websocket

	httpMux.Handle("/api/", http.StripPrefix("/api", gwMux))

	// Configure CORS options as needed
	corsEnabledMux := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace "*" with specific origin(s) for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	})

	corsHandler := corsEnabledMux.Handler(httpMux)

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *server_http_port),
		Handler: corsHandler,
	}

	logrus.Info("HTTP and websocket server started")
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Fatalln(gwServer.ListenAndServe())

}

func runGRPCServer() {

	insecureConn, err := net.Listen("tcp", fmt.Sprintf(":%d", *server_grpc_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	insecureServer := grpc.NewServer()

	// GRPC App Store Server
	nnProtoV1.RegisterNeuralNetworkServer(insecureServer, neuralnetworkserver.NewGrpcNeuralNetworkServer())

	// Register reflection service on gRPC server.
	reflection.Register(insecureServer)

	logrus.Info("GRPC server starting...")
	logrus.Fatalln(insecureServer.Serve(insecureConn))

}

func main() {

	flag.Parse()

	logrus.Info("Aayiram Karangal Neeti !")
	logrus.Info("--------->>>> Starting Neural Betwork Backend Server <<<----------")
	logrus.Infof("GRPC Port - (%d), HTTP Port - (%v)", *server_grpc_port, *server_http_port)

	go func() {
		runGRPCServer()
	}()

	runHtppServer()

}
