default: build

build:
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/neighbor.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	./bin/neighbor -filepath $(PWD)/config.yml
test:
	go test ./...

.PHONY: build install test
