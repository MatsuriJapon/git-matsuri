BUILD_DIR=./artifacts
WORK_DIR=./bin
VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-X github.com/MatsuriJapon/git-matsuri/cmd.CurrentVersion=${VERSION}"
OS ?= linux
ARCH ?= amd64
ifeq ($(OS), windows)
EXT = .exe
endif

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR} ${WORK_DIR}

.PHONY: setup
setup:
	mkdir -p ${BUILD_DIR} ${WORK_DIR}

.PHONY: lint
lint: clean
	if [ -z "$(shell which pre-commit)" ]; then pip3 install pre-commit; fi
	pre-commit install
	pre-commit run --all-files

.PHONY: test
test: clean
	go test -race -v ./...

.PHONY: build
build: clean setup
	env GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/git-matsuri-$(OS)-$(ARCH)$(EXT) .
