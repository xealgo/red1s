.PHONY: build run test vet

build:
	mkdir -p bin && go build -o ./bin/red1s ./cmd/red1s

run:
	go run ./cmd/red1s

test:
	go test -v ./...

vet: 
	go vet -v ./
