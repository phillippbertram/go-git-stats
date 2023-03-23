
dev:
	go run .

build:
	go build -o bin/ ./...

install:
	go install ./...make dev