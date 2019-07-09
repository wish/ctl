UNAME_S := $(shell uname -s | tr A-Z a-z)
SHA  	:= $(shell git rev-parse --short HEAD)
GOFILES_BUILD := $(shell find . -type f -iname "*.go")

default: bin/${UNAME_S}/ctl bin/ctl ## Builds ctl for your current operating system

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bin/linux/ctl: $(GOFILES_BUILD)
	@echo "$@"
	@GOOS=linux CGO_ENABLED=0 go build -ldflags \
	       "-X github.com/ContextLogic/ctl/cmd.Version=${SHA}" \
	       -o bin/linux/ctl github.com/ContextLogic/ctl

bin/darwin/ctl: $(GOFILES_BUILD)
	@echo "$@"
	@GOOS=darwin CGO_ENABLED=0 go build -ldflags \
		"-X github.com/ContextLogic/ctl/cmd.Version=${SHA}" \
	     -o bin/darwin/ctl github.com/ContextLogic/ctl

.PHONY: all
all: bin/linux/ctl bin/darwin/ctl ## Builds ctl binaries for linux and osx


.PHONY: clean
clean: ## Removes build artifacts
	rm -rf bin

bin/ctl: ## Make a link to the executable for this OS type for convenience
	$(shell ln -s ${UNAME_S}/ctl bin/ctl)

.PHONY: test
test: ## Runs go tests on all subdirs
	go test ./...
