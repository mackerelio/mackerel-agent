BIN = mackerel-agent
ARGS = "-conf=mackerel-agent.conf"

BUILD_FLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` "

all: clean build test

test: deps
	go test $(TESTFLAGS) github.com/mackerelio/mackerel-agent/...

build: deps
	go build -ldflags=$(BUILD_FLAGS) \
	-o build/$(BIN) \
	github.com/mackerelio/mackerel-agent

run: build
	./build/$(BIN) $(ARGS)

deps:
	go get -d github.com/mackerelio/mackerel-agent

goxc:
	go get github.com/laher/goxc

crossbuild: goxc
	goxc -build-ldflags=$(BUILD_FLAGS) \
	    -os="linux darwin windows freebsd" -arch=386 -d . \
	    -resources-include='README*' -n $(BIN)

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build run deps clean goxc crossbuild
