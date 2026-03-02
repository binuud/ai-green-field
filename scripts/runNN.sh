export FILE_STORE_PATH=${HOME}/Code/ai/ai-green-field/dataVolume
reflex $(cat .reflex) -- go run cmd/nn/main.go  -grpc_port 9090 -http_port 9080