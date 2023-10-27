build:
	@go build -o bin/gobank

run: build
	@go run main.go

test:
	@go test -v ./...
