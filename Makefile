DIRS ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)

COVERAGE_PROFILE ?= coverage.out
HTML_OUTPUT      ?= coverage.html

default: install

.PHONY: clean
clean:
	@echo "---> Cleaning"
	rm -rf ./vendor

.PHONY: enforce
enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE) -o $(HTML_OUTPUT)
	open $(HTML_OUTPUT)

.PHONY: install
install:
	@echo "---> Installing dependencies"
	dep ensure

.PHONY: lint
lint:
	@echo "---> Linting..."
	gometalinter --vendor --tests $(DIRS)

.PHONY: setup
setup:
	@echo "--> Setting up"
	go get -u -v github.com/alecthomas/gometalinter github.com/golang/dep/cmd/dep
	gometalinter --install

.PHONY: test
test:
	@echo "---> Testing"
	go test ./... -coverprofile $(COVERAGE_PROFILE)
