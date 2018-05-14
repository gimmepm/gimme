default:
	@mkdir -p bin
	@go build -o bin/gimme

test:
	@for TEST_SCRIPT in integration/test*.sh; do \
		$$TEST_SCRIPT; \
	done

clean:
	@rm -rf bin
	
.PHONY: default test clean
