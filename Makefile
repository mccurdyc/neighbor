default: build

GOCMD := $(shell which go)
GOBAK := $(GOCMD).bak

setup:
	@echo GOCMD is set to: $(GOCMD)
	@echo GOBAK is set to: $(GOBAK)
	go install github.com/golang/dep/cmd/dep
	sudo cp $(GOCMD) $(GOBAK)
	sudo cp $(PWD)/bin/go-cover $(GOCMD)

clean:
	sudo mv $(GOBAK) $(GOCMD)

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
