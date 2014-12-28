BIN = mackerel-agent
ARGS = "-conf=mackerel-agent.conf"

BUILD_FLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` "

all: clean test build

test: deps
	go test $(TESTFLAGS) ./...

build: deps
	go build -ldflags=$(BUILD_FLAGS) \
	-o build/$(BIN)

run: build
	./build/$(BIN) $(ARGS)

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get code.google.com/p/go.tools/cmd/vet
	go get github.com/laher/goxc

lint:
	go vet ./...
	golint ./...

crossbuild: deps
	$(GOPATH)/goxc -build-ldflags=$(BUILD_FLAGS) \
	    -os="linux darwin windows freebsd" -arch=386 -d . \
	    -resources-include='README*' -n $(BIN)

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build run deps clean lint crossbuild
