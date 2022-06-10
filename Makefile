install:
	go get

linter:
	golangci-lint run

test-coverage:
	CGO_ENABLED=1 go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	GOOS=darwin GOARCH=amd64 go build -o ls-lint-darwin
	GOOS=darwin GOARCH=arm64 go build -o ls-lint-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o ls-lint-linux
	GOOS=linux GOARCH=arm64 go build -o ls-lint-linux-arm64
	GOOS=windows GOARCH=amd64 go build -o ls-lint-windows.exe

build-npm:
	cp README.md npm/README.md
	make build-npm-darwin
	make build-npm-darwin-arm64
	make build-npm-linux
	make build-npm-linux-arm64
	make build-npm-windows

build-npm-darwin:
	cp ls-lint-darwin npm/bin/ls-lint-darwin
	chmod +x npm/bin/ls-lint-darwin

build-npm-darwin-arm64:
	cp ls-lint-darwin-arm64 npm/bin/ls-lint-darwin-arm64
	chmod +x npm/bin/ls-lint-darwin-arm64

build-npm-linux:
	cp ls-lint-linux npm/bin/ls-lint-linux
	chmod +x npm/bin/ls-lint-linux

build-npm-linux-arm64:
	cp ls-lint-linux-arm64 npm/bin/ls-lint-linux-arm64
	chmod +x npm/bin/ls-lint-linux-arm64

build-npm-windows:
	cp ls-lint-windows.exe npm/bin/ls-lint-windows.exe
	chmod +x npm/bin/ls-lint-windows.exe

docker-build:
	docker build -f docker/Dockerfile -t ls-lint-dev:latest .

docker-run:
	docker run --rm -v ${PWD}:/data ls-lint-dev:latest