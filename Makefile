PROJECT_NAME := "uexky"
PKG := "gitlab.com/abyss.club/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all test mod tool gen build clean

all: mod build

lint: ## Lint the files
	@golangci-lint run ./...

test: ## Run unittests
	@go test -short ${PKG_LIST}

mod: ## Get the dependencies
	@go mod tidy

tool: ## Get/Update tools
	@go get -u github.com/99designs/gqlgen
	@go get -u github.com/google/wire
	@go get -u github.com/dizzyfool/genna

gen: gengql genwire genpg

gengql: ## generate all
	@gqlgen generate

genwire:
	@cd ./uexky;wire;cd -
	@cd ./server;wire;cd -

genpg:
	@genna model -c $(pguri) -o repo/generated.go -fkw --pkg repo --gopg 9

build: mod ## Build the binary file
	@mkdir -p dist
	@go build -i -v -o dist/$(PROJECT_NAME) $(PKG) 

clean: ## Remove previous build
	@rm -f dist/$(PROJECT_NAME)
