all: cri-exporter

TAG ?= $(shell git describe --tags)
LDFLAGS := -ldflags "-w -s -X 'github.com/surik/cri-exporter.Tag=$(TAG)'"

cri-exporter: 
	@echo "Building cri-exporter..."
	go build $(LDFLAGS) -o bin/cri-exporter ./cmd

docker:
	@echo "Building images..."
	@docker build --build-arg TAG=$(TAG) -f cmd/Dockerfile -t cri-exporter .

run-docker:
	@echo "Running cri-exporter in docker..."
	@docker run -p 9000:9000 -v /var/run/cri-dockerd.sock:/var/run/cri-dockerd.sock cri-exporter

install-tools:
	@echo "Install tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	@go install gotest.tools/gotestsum@v1.12.0

lint: 
	@echo "Linting..."
	golangci-lint run

test:
	@echo "Testing..."
	@gotestsum -- -timeout 30s -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total