default:
	@mkdir -p bin
	@go build -o bin/gimme

test: default
	@for TEST_SCRIPT in integration/test*.sh; do \
		$$TEST_SCRIPT; \
	done

clean:
	@rm -rf bin

install:
	@go install github.com/gimmepm/gimme
	
.PHONY: default test clean install
