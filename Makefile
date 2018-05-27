CURRENT_VERSION = $(shell git describe --exact-match --tags $$(git log -n1 --pretty='%h') 2> /dev/null || git log -n1 --pretty='%h')

default:
	@mkdir -p bin
	@go build -ldflags "-X github.com/gimmepm/gimme/cmd.version=$(CURRENT_VERSION)" -o bin/gimme

test: default
	@for TEST_SCRIPT in integration/test*.sh; do \
		$$TEST_SCRIPT; \
	done

clean:
	@rm -rf bin dist

install:
	@go install github.com/gimmepm/gimme

docker:
	@docker build -t gimmepm/gimme:$(CURRENT_VERSION) -t gimmepm/gimme:latest .

dist:
	@mkdir -p dist/gimme-$(CURRENT_VERSION)
	@go build -ldflags "-X github.com/gimmepm/gimme/cmd.version=$(CURRENT_VERSION)" -o dist/gimme-$(CURRENT_VERSION)/gimme
	@cd dist/gimme-$(CURRENT_VERSION); tar -czf gimme-$(CURRENT_VERSION).tar.gz ./gimme
	
.PHONY: default test clean install docker dist
