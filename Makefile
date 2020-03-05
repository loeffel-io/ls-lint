install:
	go get

linter:
	golangci-lint run

test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	GOOS=darwin GOARCH=amd64 go build -o ls-lint-darwin
	GOOS=linux GOARCH=amd64 go build -o ls-lint-linux

npm-darwin:
	mkdir -p npm/ls-lint-darwin/bin
	cp ls-lint-darwin npm/ls-lint-darwin/bin/ls-lint
	chmod +x npm/ls-lint-darwin/bin/ls-lint