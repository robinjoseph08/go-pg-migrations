DIRS ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)

TFLAGS ?=

COVERAGE_PROFILE ?= coverage.out
HTML_OUTPUT      ?= coverage.html

PSQL := $(shell command -v psql 2> /dev/null)

TEST_DATABASE_USER ?= go_pg_migrations_user
TEST_DATABASE_NAME ?= go_pg_migrations

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
ifdef PSQL
	dropuser --if-exists $(TEST_DATABASE_USER)
	dropdb --if-exists $(TEST_DATABASE_NAME)
	createuser --createdb $(TEST_DATABASE_USER)
	createdb -U $(TEST_DATABASE_USER) $(TEST_DATABASE_NAME)
else
	$(error Postgres should be installed)
endif

.PHONY: test
test:
	@echo "---> Testing"
	TEST_DATABASE_USER=$(TEST_DATABASE_USER) TEST_DATABASE_NAME=$(TEST_DATABASE_NAME) go test ./... -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)
