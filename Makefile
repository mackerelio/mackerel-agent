MACKEREL_AGENT_NAME ?= "mackerel-agent"
MACKEREL_API_BASE ?= "https://mackerel.io"
VERSION = 0.45.0
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
ARGS = "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS = "linux darwin freebsd windows netbsd"

BUILD_LDFLAGS = "\
	  -X main.version=$(VERSION) \
	  -X main.gitcommit=$(CURRENT_REVISION) \
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
	gotestcover -v -race -short -covermode=atomic -coverprofile=.profile.cov -parallelpackages=4 ./...

crossbuild-package:
	mkdir -p ./build-linux-386 ./build-linux-amd64
	GOOS=linux GOARCH=386 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-386/
	GOOS=linux GOARCH=amd64 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-amd64/

crossbuild-package-kcps:
	make crossbuild-package MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16

crossbuild-package-stage:
	make crossbuild-package MACKEREL_AGENT_NAME=mackerel-agent-stage MACKEREL_API_BASE=http://0.0.0.0

rpm: rpm-v1 rpm-v2
rpm-v1: crossbuild-package
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c5 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c5 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

rpm-v2: crossbuild-package
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

deb: deb-v1 deb-v2

deb-v1: crossbuild-package
	BUILD_DIRECTORY=build-linux-386 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

deb-v2: crossbuild-package
	BUILD_DIRECTORY=build-linux-amd64 BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-kcps: rpm-kcps-v1 rpm-kcps-v2
rpm-kcps-v1: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c5 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c5 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec

rpm-kcps-v2: crossbuild-package-kcps
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec

deb-kcps: deb-kcps-v1 deb-kcps-v2

deb-kcps-v1: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps BUILD_DIRECTORY=build-linux-386 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

deb-kcps-v2: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps BUILD_SYSTEMD=1 BUILD_DIRECTORY=build-linux-amd64 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

rpm-stage: rpm-stage-v1 rpm-stage-v2
rpm-stage-v1: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c5 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" \
	-bb packaging/rpm-build/mackerel-agent-stage.spec

rpm-stage-v2: crossbuild-package-stage
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild astj/mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" \
	-bb packaging/rpm-build/mackerel-agent-stage.spec

deb-stage: deb-stage-v1 deb-stage-v2

deb-stage-v1: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage BUILD_DIRECTORY=build-linux-386 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

deb-stage-v2: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage BUILD_SYSTEMD=1 BUILD_DIRECTORY=build-linux-amd64 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -uc -us

tgz_dir = "build/tgz/$(MACKEREL_AGENT_NAME)"
tgz:
	GOOS=linux GOARCH=386 make build
	rm -rf $(tgz_dir)
	mkdir -p $(tgz_dir)
	cp mackerel-agent.sample.conf $(tgz_dir)/$(MACKEREL_AGENT_NAME).conf
	cp build/$(MACKEREL_AGENT_NAME) $(tgz_dir)/
	tar cvfz build/$(MACKEREL_AGENT_NAME)-latest.tar.gz -C build/tgz $(MACKEREL_AGENT_NAME)

check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

release: check-release-deps
	(cd _tools && cpanm -qn --installdeps .)
	perl _tools/create-release-pullrequest

commands_gen.go: commands.go
	go get github.com/motemen/go-cli/gen
	go generate

clean:
	rm -f build/$(MACKEREL_AGENT_NAME) build-linux-amd64/$(MACKEREL_AGENT_NAME) build-linux-386/$(MACKEREL_AGENT_NAME)
	go clean
	rm -f commands_gen.go

generate: commands_gen.go

.PHONY: test build run deps clean lint crossbuild cover rpm deb tgz generate crossbuild-package crossbuild-package-kcps crossbuild-package-stage rpm-v1 rpm-v2 rpm-stage rpm-stage-v1 rpm-stage-v2 rpm-kcps-v1 rpm-kcps-v2 deb-v1 deb-v2 deb-kcps deb-kcps-v1 deb-kcps-v2 deb-stage deb-stage-v1 deb-stage-v2 release check-release-deps
