SHA  := $(shell git rev-parse --short HEAD)

default: \
	build/ctl.linux \
	build/ctl.darwin

build/ctl.linux:
	@echo "$@"
	@GOOS=linux CGO_ENABLED=0 go build -ldflags \
	       "-X github.com/ContextLogic/ctl/cmd.Version=${SHA}" \
	       -o bin/linux/ctl github.com/ContextLogic/ctl

build/ctl.darwin:
	@echo "$@"
	@GOOS=darwin CGO_ENABLED=0 go build -ldflags \
		"-X github.com/ContextLogic/ctl/cmd.Version=${SHA}" \
	     -o bin/darwin/ctl github.com/ContextLogic/ctl
