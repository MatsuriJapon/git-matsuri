VERSION  := $(shell cat VERSION)
LDFLAGS  := -ldflags "-w -s -X github.com/MatsuriJapon/git-matsuri/cmd.currentVersion=${VERSION}"
BIN_DIR  := bin
BIN_NAME := git-matsuri

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: setup
setup:
	mkdir -p $(BIN_DIR)

.PHONY: lint
lint:
	if [ -z "$(shell which pre-commit)" ]; then pip3 install pre-commit; fi
	pre-commit install
	pre-commit run --all-files

.PHONY: test
test: clean
	go test -race -v ./...

.PHONY: verify
verify:
	go mod download
	go mod verify

.PHONY: build
build: clean setup
	env CGO_ENABLED=0 go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_NAME) .
