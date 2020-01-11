GO_VERSION_MIN=1.13.0
GO_CMD?=go
CGO_ENABLED?=0

# Inpired by hashicorp/vault
# https://github.com/hashicorp/vault/blob/master/Makefile
TEST?=$$($(GO_CMD) list ./...)
GOFMT_FILES?=$$(find . -name '*.go')

build:
	$(GO_CMD) fmt ./...
	$(GO_CMD) build -race -o bin/neighbor *.go

test:
	$(GO_CMD) test -v $(TEST)

.PHONY: build test
