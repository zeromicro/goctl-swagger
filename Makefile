BIN := goctl-swagger
PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = v$(shell cat version)
COMMIT_SHA ?= $(shell git describe --always)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"
GOINSTALL=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go install -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"
GOBIN ?= $(shell go env GOPATH)/bin

MAIN_ROOT ?= .

.PHONY:echo
echo:
	@echo "PKG:${PKG}"
	@echo "VERSION:${VERSION}"
	@echo "COMMIT_SHA:${COMMIT_SHA}"
	@echo "GOOS:${GOOS}"
	@echo "GOARCH:${GOARCH}"
	@echo "GOBUILD:${GOBUILD}"
	@echo "GOINSTALL:${GOINSTALL}"
	@echo "GOBIN:${GOBIN}"

.PHONY:all
all: clean build

.PHONY:install
install: download
	cd $(MAIN_ROOT) && $(GOINSTALL)

.PHONY:build
build:
	cd $(MAIN_ROOT) && $(GOBUILD) -o $(BIN)

.PHONY: test
test: build
	go test -v ./...

.PHONY: show-version
show-version:
	@echo $(VERSION)  # silent echo

.PHONY:download
download:
	go mod download -x

.PHONY:clean
clean:
	rm -rf ./$(BIN)

.PHONY:upgrade
upgrade:
	go get -u ./...

.PHONY:tidy
tidy:
	go mod tidy
