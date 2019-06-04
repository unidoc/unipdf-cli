all: build
build:
	GO111MODULE=on go build -o ./bin/unipdf ./cmd/unipdf/main.go
build-all:
	goreleaser --snapshot --skip-publish --rm-dist
release:
	goreleaser release
clean:
	rm -rf ./bin
