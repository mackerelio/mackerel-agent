MACKEREL_AGENT_NAME ?= "mackerel-agent"
MACKEREL_API_BASE ?= "https://api.mackerelio.com"
VERSION := 0.77.0
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
ARGS := "-conf=mackerel-agent.conf"
BUILD_OS_TARGETS := "linux darwin freebsd windows netbsd"
export GO111MODULE=on

BUILD_LDFLAGS := "\
	  -X main.version=$(VERSION) \
	  -X main.gitcommit=$(CURRENT_REVISION) \
	  -X github.com/mackerelio/mackerel-agent/config.agentName=$(MACKEREL_AGENT_NAME) \
	  -X github.com/mackerelio/mackerel-agent/config.apibase=$(MACKEREL_API_BASE)"

.PHONY: all
all: clean test build

.PHONY: test
test:
	go test -v -short $(TESTFLAGS) ./...

.PHONY: build
build: deps
	CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) \
	  -o build/$(MACKEREL_AGENT_NAME)

.PHONY: run
run: build
	./build/$(MACKEREL_AGENT_NAME) $(ARGS)

.PHONY: deps
deps:
	go install \
		github.com/Songmu/gocredits/cmd/gocredits \
		github.com/Songmu/goxz/cmd/goxz \

.PHONY: credits
credits: deps
	go mod tidy # not `go get` to get all the dependencies regardress of OS, architecture and build tags
	gocredits -skip-missing -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: convention
convention:
	go generate ./... && git diff --exit-code || \
	  (echo 'please `go generate ./...` and commit them' && false)

.PHONY: crossbuild
crossbuild: deps credits
	cp mackerel-agent.sample.conf mackerel-agent.conf
	goxz -build-ldflags=$(BUILD_LDFLAGS) \
		-os=linux,freebsd,netbsd -arch=386 -d ./snapshot \
		-include=mackerel-agent.conf \
		-n $(MACKEREL_AGENT_NAME) -o $(MACKEREL_AGENT_NAME)
	goxz -build-ldflags=$(BUILD_LDFLAGS) \
		-os=linux,darwin,freebsd,netbsd -arch=amd64 -d ./snapshot \
		-include=mackerel-agent.conf \
		-n $(MACKEREL_AGENT_NAME) -o $(MACKEREL_AGENT_NAME)
	goxz -build-ldflags=$(BUILD_LDFLAGS) \
		-os=linux,freebsd,netbsd -arch=arm -d ./snapshot \
		-include=mackerel-agent.conf \
		-n $(MACKEREL_AGENT_NAME) -o $(MACKEREL_AGENT_NAME)
	goxz -build-ldflags=$(BUILD_LDFLAGS) \
		-os=linux,darwin -arch=arm64 -d ./snapshot \
		-include=mackerel-agent.conf \
		-n $(MACKEREL_AGENT_NAME) -o $(MACKEREL_AGENT_NAME)
	goxz -build-ldflags=$(BUILD_LDFLAGS) \
		-os=linux -arch=mips -d ./snapshot \
		-include=mackerel-agent.conf \
		-n $(MACKEREL_AGENT_NAME) -o $(MACKEREL_AGENT_NAME)

.PHONY: cover
cover: deps
	go test -race -short -covermode=atomic -coverprofile=.profile.cov ./...

# Depending to deps looks like not needed.
# However `make build` in recipe depends deps,
# and it would install some tools as GOARCH=386 if tools are not installed.
# We should be installed tools of native architecture.
.PHONY: crossbuild-package
crossbuild-package: deps
	mkdir -p ./build-linux-386 ./build-linux-amd64 ./build-linux-arm64 ./build-linux-mips ./build-linux-armhf
	GOOS=linux GOARCH=386 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-386/
	GOOS=linux GOARCH=amd64 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-amd64/
	GOOS=linux GOARCH=arm64 make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-arm64/
	GOOS=linux GOARCH=mips make build
	mv build/$(MACKEREL_AGENT_NAME) build-linux-mips/
	GOOS=linux GOARCH=arm GOARM=6 make build # specify ARMv6 for supporting Raspberry Pi 1 / Zero
	mv build/$(MACKEREL_AGENT_NAME) build-linux-armhf/

.PHONY: crossbuild-package-kcps
crossbuild-package-kcps:
	make crossbuild-package MACKEREL_AGENT_NAME=mackerel-agent-kcps MACKEREL_API_BASE=http://198.18.0.16

.PHONY: crossbuild-package-stage
crossbuild-package-stage:
	make crossbuild-package MACKEREL_AGENT_NAME=mackerel-agent-stage MACKEREL_API_BASE=http://0.0.0.0

.PHONY: rpm
rpm: rpm-v1 rpm-v2

.PHONY: rpm-v1
rpm-v1: crossbuild-package
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" --target noarch \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64 \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

.PHONY: rpm-v2
rpm-v2: rpm-v2-x86 rpm-v2-arm

.PHONY: rpm-v2-x86
rpm-v2-x86: crossbuild-package
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64  --define "dist .el7.centos" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64 --define "dist .amzn2" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

.PHONY: rpm-v2-arm
rpm-v2-arm: crossbuild-package
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-arm64" \
	--define "_version ${VERSION}" --define "buildarch aarch64" --target aarch64  --define "dist .el7.centos" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-arm64" \
	--define "_version ${VERSION}" --define "buildarch aarch64" --target aarch64 --define "dist .amzn2" \
	-bb packaging/rpm-build/$(MACKEREL_AGENT_NAME).spec

.PHONY: deb
deb: crossbuild-package
	BUILD_DIRECTORY=build-linux-amd64 BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us
	BUILD_DIRECTORY=build-linux-arm64 BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us -aarm64
	BUILD_DIRECTORY=build-linux-mips BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us -amips
	BUILD_DIRECTORY=build-linux-armhf BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=$(MACKEREL_AGENT_NAME) _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us -aarmhf

.PHONY: rpm-kcps
rpm-kcps: rpm-kcps-v1 rpm-kcps-v2

.PHONY: rpm-kcps-v1
rpm-kcps-v1: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" --target noarch \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec
	MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64 \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec

.PHONY: rpm-kcps-v2
rpm-kcps-v2: crossbuild-package-kcps
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=mackerel-agent-kcps _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64 --define "dist .el7.centos" \
	-bb packaging/rpm-build/mackerel-agent-kcps.spec

.PHONY: deb-kcps
deb-kcps: crossbuild-package-kcps
	MACKEREL_AGENT_NAME=mackerel-agent-kcps BUILD_SYSTEMD=1 BUILD_DIRECTORY=build-linux-amd64 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: rpm-stage
rpm-stage: rpm-stage-v1 rpm-stage-v2

.PHONY: rpm-stage-v1
rpm-stage-v1: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-386" \
	--define "_version ${VERSION}" --define "buildarch noarch" --target noarch \
	-bb packaging/rpm-build/mackerel-agent-stage.spec

.PHONY: rpm-stage-v2
rpm-stage-v2: crossbuild-package-stage
	BUILD_SYSTEMD=1 MACKEREL_AGENT_NAME=mackerel-agent-stage _tools/packaging/prepare-rpm-build.sh
	docker run --rm -v "$(PWD)":/workspace -v "$(PWD)/rpmbuild":/rpmbuild mackerel/docker-mackerel-rpm-builder:c7 \
	--define "_sourcedir /workspace/packaging/rpm-build/src" --define "_builddir /workspace/build-linux-amd64" \
	--define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64 --define "dist .el7.centos" \
	-bb packaging/rpm-build/mackerel-agent-stage.spec

.PHONY: deb-stage
deb-stage: crossbuild-package-stage
	MACKEREL_AGENT_NAME=mackerel-agent-stage BUILD_SYSTEMD=1 BUILD_DIRECTORY=build-linux-amd64 _tools/packaging/prepare-deb-build.sh
	cd packaging/deb-build && debuild --no-tgz-check -rfakeroot -uc -us

tgz_dir = "build/tgz/$(MACKEREL_AGENT_NAME)"
.PHONY: tgz
tgz: credits
	GOOS=linux GOARCH=386 make build
	rm -rf $(tgz_dir)
	mkdir -p $(tgz_dir)
	cp mackerel-agent.sample.conf $(tgz_dir)/$(MACKEREL_AGENT_NAME).conf
	cp build/$(MACKEREL_AGENT_NAME) LICENSE CREDITS $(tgz_dir)/
	tar cvfz build/$(MACKEREL_AGENT_NAME)-latest.tar.gz -C build/tgz $(MACKEREL_AGENT_NAME)

.PHONY: clean
clean:
	rm -f build/$(MACKEREL_AGENT_NAME) build-linux-{386,amd64,arm64,mips,armhf}/$(MACKEREL_AGENT_NAME) CREDITS
	go clean

.PHONY: update
update:
	go get -u ./...
	go mod tidy
