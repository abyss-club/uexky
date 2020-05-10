PROJECT_NAME := "uexky"
PKG := "gitlab.com/abyss.club/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all test mod tool gen build clean

all: mod build

# lint: ## Lint the files
#   	@golint -set_exit_status ${PKG_LIST}

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
	@cd ./wire;wire;cd -

genpg:
	@genna model -c $(pguri) -o repo/generated.go -fkw --pkg repo --gopg 9

build: dep ## Build the binary file
	@go build -i -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)
