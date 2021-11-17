PROJECT_NAME=$(shell basename "$(PWD)")
VERSION=1.0.0
# go related variables.

GO_BASE=$(shell pwd)
GO_PATH=$(GO_BASE)/vendor:$(GO_BASE)
GO_BIN=$(GO_BASE)/bin
GO_FILES=$(wildcard *.go)

#default:
#	@echo "Usage of make: [build | linux | clean]"

build:
	@echo ">> go build..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=$(VERSION) -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD`" -o $(PROJECT_NAME) -v

linux:
	@echo ">> go build..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=$(VERSION) -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD`" -o $(PROJECT_NAME) -v

clean:
	@echo ">> remove..."
	rm -rf $(PROJECT_NAME)
