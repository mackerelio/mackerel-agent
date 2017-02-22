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

crossbuild-package:
	mkdir -p ./build-linux-386 ./build-linux-amd64
	GOOS=linux GOARCH=386 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-386/
	GOOS=linux GOARCH=amd64 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-amd64/

crossbuild-package-kcps:
	make crossbuild-package MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16

crossbuild-package-stage:
	mkdir -p ./build-linux-386
	GOOS=linux GOARCH=386 make build MACKEREL_AGENT_NAME=mackerel-agent-stage MACKEREL_API_BASE=http://0.0.0.0
	mv build/mackerel-agent-stage build-linux-386/

rpm: crossbuild-package
	mkdir -p rpmbuild
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-agent-packager-beta:rpm-centos6 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-agent-packager-beta:rpm-centos6 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

# TODO migrate to rpm
rpm-systemd: crossbuild-package
	mkdir -p rpmbuild
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-agent-packager-beta:rpm-centos7 \
		--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
		--define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" \
		-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME)-systemd.spec

deb: crossbuild-package
	BUILD_DIRECTORY=build-linux-386 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-kcps: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build-linux-386" \
			--define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
			-bb packaging/rpm-build/mackerel-agent-kcps.spec
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build-linux-amd64" \
			--define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" \
			-bb packaging/rpm-build/mackerel-agent-kcps.spec

deb-kcps: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps BUILD_DIRECTORY=build-linux-386 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-stage: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	rpmbuild --define "_sourcedir `pwd`/packaging/rpm-build/src" --define "_builddir `pwd`/build-linux-386" \
	      --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" \
				-bb packaging/rpm-build/mackerel-agent-stage.spec

deb-stage: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage BUILD_DIRECTORY=build-linux-386 _tools/packaging/prepare-deb-build.sh
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
	rm -f build/$(MACKEREL_AGENT_NAME) build-linux-amd64/$(MACKEREL_AGENT_NAME) build-linux-386/$(MACKEREL_AGENT_NAME)
	go clean
	rm -f commands_gen.go

generate: commands_gen.go

.PHONY: test build run deps clean lint crossbuild cover rpm deb tgz generate crossbuild-package crossbuild-package-kcps crossbuild-package-stage rpm-systemd
