BIN_DIR ?= ./bin
GO_TOOLS := \
	github.com/git-chglog/git-chglog/cmd/git-chglog \
	github.com/mattn/goveralls \

COVERAGE_PROFILE ?= coverage.out
HTML_OUTPUT      ?= coverage.html

PSQL := $(shell command -v psql 2> /dev/null)

TEST_DATABASE_USER ?= go_pg_migrations_user
TEST_DATABASE_NAME ?= go_pg_migrations

default: install

.PHONY: clean
clean:
	@echo "---> Cleaning"
	go clean
	rm -rf $(BIN_DIR) $(COVERAGE_PROFILE) $(HTML_OUTPUT)

coveralls:
	@echo "---> Sending coverage info to Coveralls"
	$(BIN_DIR)/goveralls -coverprofile=$(COVERAGE_PROFILE) -service=github

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
	go mod download

.PHONY: lint
lint: $(BIN_DIR)/golangci-lint
	@echo "---> Linting"
	$(BIN_DIR)/golangci-lint run

.PHONY: postgres
postgres:
	@echo "---> Starting Postgres"
	docker run \
		--name go-pg-postgres \
		--rm \
		-e POSTGRES_DB=$(TEST_DATABASE_NAME) \
		-e POSTGRES_HOST_AUTH_METHOD=trust \
		-e POSTGRES_USER=$(TEST_DATABASE_USER) \
		-p 5432:5432 \
		postgres:11

.PHONY: release
release:
	@echo "---> Creating new release"
ifndef tag
	$(error tag must be specified)
endif
	$(BIN_DIR)/git-chglog --output CHANGELOG.md --next-tag $(tag)
	sed -i "" "s/version-.*-green/version-$(tag)-green/" README.md
	git add CHANGELOG.md README.md
	git commit -m $(tag)
	git tag $(tag)
	git push origin master --tags

.PHONY: setup
setup: $(BIN_DIR)/golangci-lint
	@echo "--> Setting up"
	GOBIN=$(PWD)/$(BIN_DIR) go install $(GO_TOOLS)
ifdef PSQL
	dropdb --if-exists $(TEST_DATABASE_NAME) 2> /dev/null && \
		dropuser --if-exists $(TEST_DATABASE_USER) && \
		createuser --createdb $(TEST_DATABASE_USER) && \
		createdb -U $(TEST_DATABASE_USER) $(TEST_DATABASE_NAME) || true
else
	$(warning Postgres not installed locally; run `make postgres` to start within Docker)
endif

$(BIN_DIR)/golangci-lint:
	@echo "---> Installing linter"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.61.0

.PHONY: test
test:
	@echo "---> Testing"
	TEST_DATABASE_USER=$(TEST_DATABASE_USER) TEST_DATABASE_NAME=$(TEST_DATABASE_NAME) go test . -coverprofile $(COVERAGE_PROFILE)
