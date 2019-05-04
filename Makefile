deps:
	go get golang.org/x/lint/golint

lint: deps
	test -z "$(gofmt -s -l . | grep -v '^vendor' | tee /dev/stderr)"
	golint -set_exit_status ./...
	go build ./...
	go test -race -v ./...
	go vet ./...

build: lint
	env GOOS=windows GOARCH=amd64 go build -o ./artifacts/windows_amd64/git-matsuri.exe .
	env GOOS=linux GOARCH=amd64 go build -o ./artifacts/linux_amd64/git-matsuri .
	env GOOS=darwin GOARCH=amd64 go build -o ./artifacts/darwin_amd64/git-matsuri .

.PHONY: deps lint build
