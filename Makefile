BUILD_DIR=./artifacts
clean:
	rm -rf ${BUILD_DIR}

deps:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.16.0

lint: deps
	golangci-lint run --new-from-rev=HEAD~
	go build ./...
	go test -race -v ./...

build: clean lint
	mkdir -p ${BUILD_DIR}
	env GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/windows_amd64/git-matsuri.exe .
	env GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/linux_amd64/git-matsuri .
	env GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/darwin_amd64/git-matsuri .

.PHONY: deps lint build clean
