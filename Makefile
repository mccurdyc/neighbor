default: build

setup:
	./build/setup.sh

build:
	go get -u -v github.com/golang/dep/cmd/dep
	dep ensure -v
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/main.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	./bin/neighbor -filepath $(PWD)/config.json

test:
	go test ./...

.PHONY: setup build install run test
