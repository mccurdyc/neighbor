build:
	go fmt ./...
	go build -race -o bin/neighbor cmd/neighbor/main.go

test:
	go test -v ./...

.PHONY: build test
