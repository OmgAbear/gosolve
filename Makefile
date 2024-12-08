.PHONY: build run test clean
## Variables
GO = go
BINARY_NAME = number-index-service

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:## Build the go application
	$(GO) mod tidy
	$(GO) build -o $(BINARY_NAME) cmd/server/main.go
run: build ## Run the go server
	NUMBERS_FILE_LOCATION=static/input.txt ./$(BINARY_NAME)
test:## Run go tests
	$(GO) test ./... -v
clean:## Clean go build artifacts
	$(GO) clean
	rm -f $(BINARY_NAME)
install:## Install go dependencies
	$(GO) mod download
docker-run:## Run docker services
	cd dev
	docker-compose -f dev/docker-compose.yaml up --build -d
docker-stop:## Stop docker services
	cd dev
	docker-compose -f dev/docker-compose.yaml down
