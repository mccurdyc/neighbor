default: build

build:
	go fmt cmd/neighbor
	go build -o bin/neighbor cmd/neighbor

install: build
	cp bin/neighbor /usr/local/bin

test:
	go test cmd/neighbor

.PHONY: build install test
