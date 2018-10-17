default: build

setup:
	go get -u -v github.com/golang/dep/cmd/dep
	dep ensure -v
	./build/setup.sh

clean:
	./build/clean.sh

build:
	go fmt ./...
	go build -o bin/neighbor cmd/neighbor/main.go

install: build
	cp bin/neighbor /usr/local/bin

run: build
	COVERPROFILE_FNAME="neighbor-coverprofile.out" ./bin/neighbor -filepath $(PWD)/config.json

test:
	go test ./...

.PHONY: build install setup run test clean
