BUILD_DIR=./artifacts
WORK_DIR=./bin
VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-X github.com/MatsuriJapon/git-matsuri/cmd.CurrentVersion=${VERSION}"
OS ?= linux
ARCH ?= amd64
ifeq ($(OS), windows)
EXT = .exe
endif

clean:
	rm -rf ${BUILD_DIR} ${WORK_DIR}

setup:
	mkdir -p ${BUILD_DIR} ${WORK_DIR}

lint: clean setup
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b ${WORK_DIR} v1.27.0
	${WORK_DIR}/golangci-lint run

test: clean
	go test -race -v ./...

build: clean setup
	env GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/git-matsuri-$(OS)-$(ARCH)$(EXT) .

.PHONY: setup test lint build clean
