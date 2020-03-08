.PHONY: all build deps image lint migrate test race msan vet
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	go build -ldflags "-X github.com/delivc/identity/cmd.Version=`git rev-parse HEAD`"

deps: ## Install dependencies.
	@go get -v -d ./...
	@go get -u github.com/gobuffalo/pop/soda
	@go get -u golang.org/x/lint
	@go mod download

image: ## Build the Docker image.
	docker build .

lint: ## Lint the code.
	@golint -set_exit_status $(CHECK_FILES)

migrate_dev: ## Run database migrations for development.
	hack/migrate.sh development

migrate_test: ## Run database migrations for test.
	hack/migrate.sh test

test: ## Run tests.
	go test -short -p 1 -v $(CHECK_FILES)

race:
	go test -race -short ${CHECK_FILES}

msan:
	@go test -msan -short ${CHECK_FILES}

coverage:
	./hack/coverage.sh

coverhtml:
	./hack/coverage.sh html

vet: # Vet the code
	go vet $(CHECK_FILES)