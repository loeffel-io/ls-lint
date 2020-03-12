install:
	go get

linter:
	golangci-lint run

test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	GOOS=darwin GOARCH=amd64 go build -o ls-lint-darwin
	GOOS=linux GOARCH=amd64 go build -o ls-lint-linux
	GOOS=windows GOARCH=amd64 go build -o ls-lint-windows.exe

build-npm:
	cp ls-lint-darwin npm/ls-lint/bin/ls-lint-darwin
	cp ls-lint-linux npm/ls-lint/bin/ls-lint-linux
	cp ls-lint-windows.exe npm/ls-lint/bin/ls-lint-windows.exe
	chmod +x npm/ls-lint/bin/ls-lint-darwin
	chmod +x npm/ls-lint/bin/ls-lint-linux
	chmod +x npm/ls-lint/bin/ls-lint-windows.exe
