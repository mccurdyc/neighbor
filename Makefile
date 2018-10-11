default: build

setup:
	go install github.com/golang/dep/cmd/dep
	sudo ./build/setup.sh

clean:
	sudo ./build/clean.sh

build:
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/main.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	./bin/neighbor -filepath $(PWD)/config.yml

test:
	go test ./...

.PHONY: build install setup run test clean
