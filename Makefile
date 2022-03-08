.DEFAULT_GOAL := build

export GO111MODULE=on
export CGO_ENABLED=0
export VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
BINARY=vuemonit
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -s -w"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: front
front:
	cd front && npm install && ./node_modules/.bin/quasar build

.PHONY: build
build: ## Build
	go build $(LDFLAGS) -o $(BINARY) 

.PHONY: fullbuild
fullbuild: proto front
	mv front/dist/spa/index.html front/dist/spa/main.html
	statik -src=./front/dist/spa/ -f
	mv front/dist/spa/main.html front/dist/spa/index.html
	go build $(LDFLAGS) -o $(BINARY)

.PHONY: compressed
compressed: fullbuild
	upx --best $(BINARY)

.PHONY: proto
proto: ## Generate protobuf
	protoc --go_out=. implem/storage.storm/*.proto

.PHONY: docker
docker: ## Build the docker image
	docker build -t $(BINARY):latest -t $(BINARY):$(BUILD) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

.PHONY: release
release: ## Create a new release on Github
	goreleaser

.PHONY: snapshot
snapshot: ## Create a new snapshot release
	goreleaser --snapshot --rm-dist

.PHONY: lint
lint: ## Runs the linter
	$(GOPATH)/bin/golangci-lint run --exclude-use-default=false

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test -race -coverprofile="coverage.txt" ./...

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi
