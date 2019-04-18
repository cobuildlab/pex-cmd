all: build

clean:
	rm -f pex-cmd

prepare:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

build: prepare clean
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

install: prepare clean
	CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo

run: prepare
	go run main.go

.PHONY: install prepare build clean run

