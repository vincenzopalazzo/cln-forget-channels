CC=go
FMT=gofmt
NAME=cln4go-plugin
BASE_DIR=/script
OS=linux
ARCH=386
ARM=
GORPC_COMMIT=4471a927bb9937a45a9ece876c3e00f093727fc3

default: fmt lint
	$(CC) build -o $(NAME) cmd/plugin.go

fmt:
	$(CC) fmt ./...

lint:
	golangci-lint run

check:
	$(CC) test -v ./...

build:
	env GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARM) $(CC) build -o $(NAME)-$(OS)-$(ARCH) cmd/plugin.go

dep:
	$(CC) mod vendor

force:
	$(CC) get -u all
	$(CC) get -u github.com/vincenzopalazzo/cln4go@$(GORPC_COMMIT)
	@make dep

integration:
	@make default
	cd cln-integration; make fmt
	cd cln-integration; make check


clean:
	cd cln-integration; make clean
