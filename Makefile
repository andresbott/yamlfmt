SHELL := /bin/bash


# ======================================================================================
default: help;

fmt: ## format go code and run mod tidy
	@go fmt ./...
	@go mod tidy

.PHONY: test
test: ## run go tests
	@go test ./...  -cover

lint: ## run go linter
	@golangci-lint run

benchmark: ## run go benchmarks
	@go test -run=^$$ -bench=. ./...

verify: fmt test benchmark lint ## run all verification and code structure tiers

build: ## builds a snapshot build using goreleaser
	@goreleaser --snapshot --rm-dist

release: verify ## release a new version of goback
	@:$(call check_defined, version, "version defined: call with version=\"v1.2.3\"")
	@git diff --quiet || ( echo 'git is in dirty state' ; exit 1 )
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )
	@git tag -d $(version) || true # delete tag if it exists, allows to overwrite tags
	@git push --delete origin $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push origin $(version)
	@goreleaser --rm-dist


help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m·%-20s\033[0m %s\n", $$1, $$2}'