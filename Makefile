BINARY_NAME=n2h
GOOS=darwin
GOARCH=arm64
DBNAME=$(DATABASE_NAME)
LDFLAGS=-ldflags "-X=main.PWD=$(PWD)"

build:
	GOARCH=${GOARCH} GOOS=${GOOS} go build -o ${BINARY_NAME} ${LDFLAGS} main.go

run:
	./${BINARY_NAME}

install:
	go install ${LDFLAGS}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}