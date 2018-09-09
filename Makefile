default: build

build:
	go fmt cmd/neighbor
	GO111MODULE=on go build -mod=vendor -o bin/neighbor cmd/neighbor

install: build
	cp bin/neighbor /usr/local/bin

test:
	GO111MODULE=on go test cmd/neighbor

.PHONY: build install test
