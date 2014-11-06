BIN = mackerel-agent
ARGS = "-conf=mackerel-agent.conf"

all: clean build test

test: deps
	go test $(TESTFLAGS) github.com/mackerelio/mackerel-agent/...

build: deps
	go build \
	-ldflags="\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` " \
	-o build/$(BIN) \
	github.com/mackerelio/mackerel-agent

run: build
	./build/$(BIN) $(ARGS)

deps:
	go get -d github.com/mackerelio/mackerel-agent

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build run deps clean
