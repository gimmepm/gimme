CURRENT_VERSION = $(shell git describe --exact-match --tags $$(git log -n1 --pretty='%h') 2> /dev/null || git log -n1 --pretty='%h')

default:
	@mkdir -p bin
	@go build -ldflags "-X github.com/gimmepm/gimme/cmd.version=$(CURRENT_VERSION)" -o bin/gimme

test: default
	@for TEST_SCRIPT in integration/test*.sh; do \
		$$TEST_SCRIPT; \
	done

clean:
	@rm -rf bin

install:
	@go install github.com/gimmepm/gimme

docker:
	@docker build -t gimmepm/gimme:$(CURRENT_VERSION) .
	
.PHONY: default test clean install docker
