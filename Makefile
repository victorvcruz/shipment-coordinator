APP_NAME=shipment-coordinator
BUILD_DIR=bin
MAIN_PACKAGE=./cmd

test:
	go test ./... -v

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PACKAGE)

up:
	docker-compose up --build

down:
	docker-compose down

clean:
	rm -rf $(BUILD_DIR)

lint:
	golangci-lint run -c .golangci.yaml

run:
	go run $(MAIN_PACKAGE)/main.go