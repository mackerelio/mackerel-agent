BIN = mackerel-agent
ARGS = "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS = "linux darwin freebsd windows netbsd"

BUILD_LDFLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT=`git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION=`git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` \
	  -X github.com/mackerelio/mackerel-agent/config.agentName=$(MACKEREL_AGENT_NAME) \
	  -X github.com/mackerelio/mackerel-agent/config.apibase=$(MACKEREL_API_BASE)"

all: clean test build

test: lint
	go test $(TESTFLAGS) ./...

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) \
	-o build/$(BIN)

run: build
	./build/$(BIN) $(ARGS)

deps:
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
	    -os="linux darwin freebsd netbsd" -arch="386 amd64 arm" -d . -n $(BIN)

cover: deps
	gotestcover -v -short -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

rpm:
	GOOS=linux GOARCH=386 make build
	cp mackerel-agent.sample.conf packaging/rpm/src/mackerel-agent.conf
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm/src" --define "_builddir `pwd`/build" -ba packaging/rpm/mackerel-agent.spec

deb:
	GOOS=linux GOARCH=386 make build
	cp build/$(BIN)        packaging/deb/debian/mackerel-agent.bin
	cp mackerel-agent.sample.conf packaging/deb/debian/mackerel-agent.conf
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

release:
	_tools/releng

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build run deps clean lint crossbuild cover rpm deb
