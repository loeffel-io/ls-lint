install:
	go get

linter:
	golangci-lint run

test:
	make linter

build:
	GOOS=darwin GOARCH=amd64 go build -o ls-lint-darwin
	GOOS=linux GOARCH=amd64 go build -o ls-lint-linux
	chmod +x ls-lint-darwin
	chmod +x ls-lint-linux