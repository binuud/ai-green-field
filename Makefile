# Maintainer Binu Udayakumar <binu@dronasys.com>

REPO=dronasys

MULTI_PLATFORM_DOCKER = --platform=linux/amd64,linux/arm64

# Optional environmental variables


BASE_DIR	:=	$(shell git rev-parse --show-toplevel)
BIN			:=	$(BASE_DIR)/bin

## Needs protoc to be installed

ifndef PROTOC
PROTOC = protoc
endif

ifndef PROTOWEB
PROTOWEB = protoc-gen-grpc-web
endif

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help vendor vendor-update
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

tidy: ## run tidy
	go mod tidy

vendor: ## run vendor
	go mod tidy
	go mod vendor

protos: ## generate protos for neural network
	$(PROTOC) -I=./proto/.  \
	--go_out=./gen/go/ \
	--go_opt paths=source_relative \
	--go-grpc_out=./gen/go/ \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out=./gen/go/ \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--grpc-gateway-ts_out=./gen/go/ \
	--grpc-gateway-ts_opt paths=source_relative \
	--grpc-gateway-ts_opt generate_unbound_methods=true \
	--oas_out ./gen/go/v1/neuralNetwork/ \
	proto/v1/neuralNetwork/neuralNetwork.proto	\
	proto/v1/neuralNetwork/nnService.proto 
	yq eval ./gen/go/v1/neuralNetwork/openapi.yaml -o=json -P > ./gen/go/v1/neuralNetwork/openapi.json

run-staging: ## run staging test
	./scripts/runStaging.sh	

test: ## Run tests
	go test -v ./pkg/neuralNetwork/...
