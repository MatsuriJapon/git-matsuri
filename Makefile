BUILD_DIR=./artifacts
TMP_BIN_DIR=./bin
VERSION := $(shell cat VERSION)
LDFLAGS=-ldflags "-X github.com/MatsuriJapon/git-matsuri/cmd.CurrentVersion=${VERSION}"
clean:
	rm -rf ${BUILD_DIR} ${TMP_BIN_DIR}

setup:
	mkdir -p ${BUILD_DIR}
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.23.3

lint: setup
	${TMP_BIN_DIR}/golangci-lint run --new-from-rev=HEAD~
	go build ./...
	go test -race -v ./...

build: clean lint
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/windows_amd64/git-matsuri.exe .
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/linux_amd64/git-matsuri .
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/darwin_amd64/git-matsuri .

.PHONY: deps lint build clean
