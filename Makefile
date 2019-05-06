BUILD_DIR=./artifacts
clean:
	rm -rf ${BUILD_DIR}

deps:
	go get golang.org/x/lint/golint

lint: deps
	test -z "$(gofmt -s -l . | grep -v '^vendor' | tee /dev/stderr)"
	golint -set_exit_status ./...
	go build ./...
	go test -race -v ./...
	go vet ./...

build: clean lint
	mkdir -p ${BUILD_DIR}
	env GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/windows_amd64/git-matsuri.exe .
	env GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/linux_amd64/git-matsuri .
	env GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/darwin_amd64/git-matsuri .

.PHONY: deps lint build clean
