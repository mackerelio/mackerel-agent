BIN = mackerel-agent
ARGS = "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS = "linux darwin freebsd windows"

BUILD_LDFLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` \
	  -X github.com/mackerelio/mackerel-agent/config.apibase \"$(MACKEREL_API_BASE)\""

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
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/vet
	go get golang.org/x/tools/cmd/cover
	go get github.com/laher/goxc
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: deps
	go vet ./...
	rm -f $(LINT_RET)
	for os in "$(BUILD_OS_TARGETS)"; do \
		if [ $$os != "windows" ]; then \
			GOOS=$$os golint ./... | grep -v '_string.go:' | tee -a $(LINT_RET); \
		else \
			GOOS=$$os golint --min_confidence=0.9 ./... | grep -v '_string.go:' | tee -a $(LINT_RET); \
		fi \
	done
	test ! -s $(LINT_RET)

crossbuild: deps
	cp mackerel-agent.sample.conf mackerel-agent.conf
	goxc -build-ldflags=$(BUILD_LDFLAGS) \
	    -os=$(BUILD_OS_TARGETS) -arch="386 amd64 arm" -d . \
	    -resources-include='README*,mackerel-agent.conf' -n $(BIN) \
	    -main-dirs-exclude ./wix

cover: deps
	tool/cover.sh

rpm: mkr
	GOOS=linux GOARCH=386 make build
	cp mackerel-agent.sample.conf packaging/rpm/src/mackerel-agent.conf
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm/src" --define "_builddir `pwd`/build" -ba packaging/rpm/mackerel-agent.spec

deb: mkr
	GOOS=linux GOARCH=386 make build
	cp build/$(BIN)               packaging/deb/debian/mackerel-agent.bin
	cp build/mkr                  packaging/deb/debian/mkr.bin
	cp mackerel-agent.sample.conf packaging/deb/debian/mackerel-agent.conf
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

mkr:
	mkdir -p build
	curl -Ls -o build/mkr https://github.com/mackerelio/mkr/releases/download/v0.3.1/mkr_linux_386
	chmod +x build/mkr

release:
	tool/releng

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build run deps clean lint crossbuild cover rpm deb mkr
