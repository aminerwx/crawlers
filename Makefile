BINARY_NAME=crawlers
DATABASE_URL=postgres://postgres@localhost:5432/crawlers?sslmode=disable
build:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go

run: build
	bin/${BINARY_NAME}-linux

clean:
	go clean
	rm bin/${BINARY_NAME}-linux

up:
	migrate -database ${DATABASE_URL} -path migration up

down:
	migrate -database ${DATABASE_URL} -path migration down

migrate:
	ifeq up	
		up	
	endif
	ifeq down 
		down
	endif
