SHORT?= $(shell git rev-parse --short HEAD)
VERSION?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
DATE=$(shell gdate -u +%s)
LDFLAGS = -ldflags "-s -w -X github.com/slillibri/xterrafile/cmd.version=${VERSION} -X github.com/slillibri/xterrafile/cmd.commit=${SHORT} -X github.com/slillibri/xterrafile/cmd.date=${DATE}"

clean:
	rm -fv bin/*
	rm -fv release/*

build:
	mkdir -v -p bin
	go build -o bin/xterrafile -v ${LDFLAGS}

compile:
	echo "Compiling for every OS and Platform"
	mkdir -v -p bin
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/xterrafile-linux-amd64
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/xterrafile-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/xterrafile-darwin-arm64

release: clean compile
	mkdir -v -p release
	tar zcf release/xterrafile-linux-amd64.tar.gz bin/xterrafile-linux-amd64
	tar zcf release/xterrafile-darwin-amd64.tar.gz bin/xterrafile-darwin-amd64
	tar zcf release/xterrafile-darwin-arm64.tar.gz bin/xterrafile-darwin-arm64

all: build