LDFLAGS := -s -w

VERSION ?= 1.0.0
BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)  $(shell git log -1 --pretty=%s)

build:
	@env CGO_ENABLED=0 \
	go build -trimpath \
		-ldflags "$(LDFLAGS) \
		-X 'uaProxy/bootstrap.Version=$(VERSION)' \
		-X 'uaProxy/bootstrap.BuildDate=$(BUILD_DATE)' \
		-X 'uaProxy/bootstrap.GitCommit=$(GIT_COMMIT)'" \
		.

debug:
	@CompileDaemon -build="make build" -command="./uaProxy --debug --stats"

clean:
	@rm uaProxy
	@rm -rf dist

check:
	gofumpt -l -w .
	golangci-lint run

goreleaser-check:
	goreleaser release --snapshot --clean

.PHONY: debug clean check
