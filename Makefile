.PHONY: build test fmt vet lint clean install help

BINARY := itg
PKG    := ./...

help:
	@echo "make build     build the itg binary"
	@echo "make test      run unit tests"
	@echo "make vet       run go vet"
	@echo "make fmt       gofmt -w on all .go files"
	@echo "make install   go install the binary into GOBIN"
	@echo "make clean     remove build artifacts"

build:
	go build -o $(BINARY) ./cmd/itg

test:
	go test -race -count=1 $(PKG)

vet:
	go vet $(PKG)

fmt:
	gofmt -w .

install:
	go install ./cmd/itg

clean:
	rm -f $(BINARY)
