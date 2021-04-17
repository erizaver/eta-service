LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=eta-service

.PHONY: build
build: test
	go build -v -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd

.PHONY: run
run:
	go run cmd/main.go --local-config=values/values_local.yaml

.PHONY: deps
deps:
	go mod tidy
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: docker-build
docker-build:
	docker build -t $(PROJECT_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 $(PROJECT_NAME)

.PHONY: redis-run
redis-run:
	docker pull redis
	docker run -p 6379:6379 -d redis

.PHONY: docker-compose-up
docker-compose-up:
	docker compose up

.PHONY: depsGen
protoDeps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	brew tap go-swagger/go-swagger
	brew install go-swagger

.PHONY: generate
generate:
	protoc -I ./pb/ \
	  -I . \
	  --grpc-gateway_out ./pkg \
	  --go_out ./pkg \
	  --go-grpc_out ./pkg \
      --grpc-gateway_opt logtostderr=true \
      --grpc-gateway_opt paths=source_relative \
      api/eta-service.proto

.PHONY: generateExternalClients
generateExternalClients:
	swagger generate client -f ./internal/pkg/externalclients/carservice/car-swagger.yml -t ./internal/pkg/externalclients/carservice/ --default-scheme=https
	swagger generate client -f ./internal/pkg/externalclients/predictservice/predict-swagger.yml -t ./internal/pkg/externalclients/predictservice/ --default-scheme=https

.PHONY: generateMocks
generateMocks:
	mockery --dir ./internal/pkg/etaservice --output ./internal/pkg/etaservice/mocks --all

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./...

.PHONY: test
test:
	go test ./...