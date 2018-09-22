default: build

setup:
	go install github.com/golang/dep/cmd/dep

build:
	dep ensure
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/main.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	./bin/neighbor -filepath $(PWD)/config.yml

test:
	go test ./...

.PHONY: build install setup run test
