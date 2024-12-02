SHELL := /bin/bash

# test task
#
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./tests/...