all: build

clean:
	rm -f pex-cmd

prepare:
	go get github.com/tools/godep
	godep restore
	godep save

build: prepare clean
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo

install: prepare clean
	CGO_ENABLED=0 GOOS=linux godep go install -a -installsuffix cgo

run: prepare
	godep go run main.go

.PHONY: install prepare build clean run

