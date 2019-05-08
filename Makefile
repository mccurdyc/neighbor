build:
	go fmt ./...
	GO111MODULE=on go build -race -o bin/neighbor cmd/neighbor/main.go

test:
	GO111MODULE=on go test -v ./...

.PHONY: build test
