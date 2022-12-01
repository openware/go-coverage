build: all

all: go-coverage

go-coverage:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o bin/go-coverage

clean:
	rm -rf bin/*
