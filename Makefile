default: build

setup:
ifeq ($(shell go version | grep -e 'go1.11'),)
	@echo "make sure that you have Go >= 1.11 installed"
	exit 1
endif
	export GO111MODULE=on
	go get ./...
	./build/setup.sh $(PWD)

build:
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/main.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	./bin/neighbor -file=$(PWD)/config.json

test:
	go test ./...

.PHONY: setup build install run test
