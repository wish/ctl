SHA  := $(shell git rev-parse --short HEAD)

default: \
	build/wishctl.linux \
	build/wishctl.darwin

build/wishctl.linux:
	@echo "$@"
	@GOOS=linux CGO_ENABLED=0 go build -ldflags \
	       "-X github.com/ContextLogic/wishctl/cmd.Version=${SHA}" \
	       -o bin/linux/wishctl github.com/ContextLogic/wishctl

build/wishctl.darwin:
	@echo "$@"
	@GOOS=darwin CGO_ENABLED=0 go build -ldflags \
		"-X github.com/ContextLogic/wishctl/cmd.Version=${SHA}" \
	     -o bin/darwin/wishctl github.com/ContextLogic/wishctl
