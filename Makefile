MACKEREL_AGENT_NAME ?= "mackerel-agent"
MACKEREL_API_BASE ?= "https://mackerel.io"
MACKEREL_AGENT_VERSION ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/[-+].*$$//')
ARGS = "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS = "linux darwin freebsd windows netbsd"

BUILD_LDFLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT=`git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION=$(MACKEREL_AGENT_VERSION) \
	  -X github.com/mackerelio/mackerel-agent/config.agentName=$(MACKEREL_AGENT_NAME) \
	  -X github.com/mackerelio/mackerel-agent/config.apibase=$(MACKEREL_API_BASE)"

all: clean test build

test: lint
	go test -v -short $(TESTFLAGS) ./...

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) \
	-o build/$(MACKEREL_AGENT_NAME)

run: build
	./build/$(MACKEREL_AGENT_NAME) $(ARGS)

deps: generate
	go get -d -v -t ./...
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/laher/goxc
	go get github.com/mattn/goveralls

lint: deps
	go tool vet -all .
	_tools/go-linter $(BUILD_OS_TARGETS)

crossbuild: deps
	cp mackerel-agent.sample.conf mackerel-agent.conf
	goxc -build-ldflags=$(BUILD_LDFLAGS) \
	    -os="linux darwin freebsd netbsd" -arch="386 amd64 arm" -d . -n $(MACKEREL_AGENT_NAME)

cover: deps
	gotestcover -v -short -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

rpm:
	GOOS=linux GOARCH=386 make build
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build" \
	    -ba packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

deb:
	GOOS=linux GOARCH=386 make build
	MACKEREL_AGENT_VERSION=$(MACKEREL_AGENT_VERSION) MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) \
	  _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

release:
	_tools/releng

commands_gen.go: commands.go
	go get github.com/motemen/go-cli/gen
	go generate

logging/level_string.go: logging/level.go
	go get golang.org/x/tools/cmd/stringer
	go generate ./logging

clean:
	rm -f build/$(MACKEREL_AGENT_NAME)
	go clean
	rm -f commands_gen.go

generate: commands_gen.go

.PHONY: test build run deps clean lint crossbuild cover rpm deb generate
