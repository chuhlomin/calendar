.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: vet
## vet: runs the linter
vet:
	@go vet ./...

.PHONY: test
## test: runs the tests
test:
	@go test ./...

.PHONY: run
## run: runs the binary
run:
	@go build -o grid .
	@./grid
