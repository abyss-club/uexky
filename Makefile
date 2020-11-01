PROJECT_NAME := "uexky"
PKG := "gitlab.com/abyss.club/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all test mod tool gen build clean

all: mod gen build lint

lint: ## Lint the files
	@golangci-lint run ./...

test: ## Run unittests
	@go test -short ${PKG_LIST}

mod: ## Get the dependencies
	@go mod tidy

tool: ## Get/Update tools
	@go get -u github.com/99designs/gqlgen
	@go get -u github.com/google/wire

gen: gengql genwire

gengql: ## generate all
	@gqlgen generate

genwire:
	@cd ./uexky;wire;cd -
	@cd ./auth;wire;cd -
	@cd ./server;wire;cd -

build: mod ## Build the binary file
	@mkdir -p dist
	@go build -i -v -o dist/$(PROJECT_NAME) $(PKG) 

clean: ## Remove previous build
	@rm -f dist/$(PROJECT_NAME)
