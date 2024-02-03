BINARY_NAME=crawlers

build:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go

run: build
	bin/${BINARY_NAME}-linux

clean:
	go clean
	rm bin/${BINARY_NAME}-linux
