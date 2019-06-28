all: clean build ; 
build:
	go build -o ./bin/producer ./producer/main.go
	go build -o ./bin/consumer ./consumer/main.go
clean:
	go clean && rm -rf ./bin