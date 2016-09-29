MACKEREL_AGENT_NAME ?= "mackerel-agent"
MACKEREL_API_BASE ?= "https://mackerel.io"
MACKEREL_AGENT_VERSION ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/[-+].*$$//')
ARGS = "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS = "linux darwin freebsd windows netbsd"
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')


BUILD_LDFLAGS = "\
	  -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT=`git rev-parse --short HEAD` \
	  -X github.com/mackerelio/mackerel-agent/version.VERSION=$(CURRENT_VERSION) \
	  -X github.com/mackerelio/mackerel-agent/config.agentName=$(MACKEREL_AGENT_NAME) \
	  -X github.com/mackerelio/mackerel-agent/config.apibase=$(MACKEREL_API_BASE)"

check-variables:
	echo "CURRENT_VERSION: ${CURRENT_VERSION}"
	echo "MACKEREL_AGENT_NAME: ${MACKEREL_AGENT_NAME}"
	echo "MACKEREL_AGENT_VERSION: ${MACKEREL_AGENT_VERSION}"
	echo "MACKEREL_API_BASE: ${MACKEREL_API_BASE}"

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
	go get github.com/golang/lint/golint
	go get github.com/pierrre/gotestcover
	go get github.com/laher/goxc
	go get github.com/mattn/goveralls

lint: deps
	go tool vet -all -printfuncs=Criticalf,Infof,Warningf,Debugf,Tracef .
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
	      --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
				-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	GOOS=linux GOARCH=amd64 make build
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build" \
			--define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" \
			-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

deb:
	GOOS=linux GOARCH=386 make build
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-kcps:
	make build MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16 GOOS=linux GOARCH=386
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build" \
			--define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
			-bb packaging/rpm-build/mackerel-agent-kcps.spec
	make build MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16 GOOS=linux GOARCH=amd64
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build" \
			--define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" \
			-bb packaging/rpm-build/mackerel-agent-kcps.spec

deb-kcps:
	make build MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16 GOOS=linux GOARCH=386
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-stage:
	make build MACKEREL_AGENT_NAME=mackerel-agent-stage MACKEREL_API_BASE=http://0.0.0.0 GOOS=linux GOARCH=386
	MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build" \
	      --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
				-bb packaging/rpm-build/mackerel-agent-stage.spec

deb-stage:
	make build MACKEREL_AGENT_NAME=mackerel-agent-stage MACKEREL_API_BASE=http://0.0.0.0 GOOS=linux GOARCH=386
	MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

tgz_dir = "build/tgz/$(MACKEREL_AGENT_NAME)"
tgz:
	GOOS=linux GOARCH=386 make build
	rm -rf $(tgz_dir)
	mkdir -p $(tgz_dir)
	cp mackerel-agent.sample.conf $(tgz_dir)/$(MACKEREL_AGENT_NAME).conf
	cp build/$(MACKEREL_AGENT_NAME) $(tgz_dir)/
	tar cvfz build/$(MACKEREL_AGENT_NAME)-latest.tar.gz -C build/tgz $(MACKEREL_AGENT_NAME)

release:
	_tools/releng

commands_gen.go: commands.go
	go get github.com/motemen/go-cli/gen
	go generate

clean:
	rm -f build/$(MACKEREL_AGENT_NAME)
	go clean
	rm -f commands_gen.go

generate: commands_gen.go

.PHONY: test build run deps clean lint crossbuild cover rpm deb tgz generate
