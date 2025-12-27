# ============================
# Project Config
# ============================
APP_NAME=auth_service
GO=go
MOCKGEN=mockgen

# ============================
# Paths
# ============================
EXECUTOR_SRC=auth_service/infrastructure/database/executor.go
MOCK_DIR=auth_service/infrastructure/database/mocks
MOCK_FILE=$(MOCK_DIR)/mock_executor.go

REPO_PKG=./auth_service/internal/repository
COVERAGE_FILE=coverage.out

# ============================
# Phony
# ============================
.PHONY: help mock test test-all coverage coverage-html clean

# ============================
# Help
# ============================
help:
	@echo "Available commands:"
	@echo "  make mock           Generate gomock mocks"
	@echo "  make test           Run repository unit tests"
	@echo "  make test-all       Run all unit tests"
	@echo "  make coverage       Run tests with coverage"
	@echo "  make coverage-html  Open coverage in browser"
	@echo "  make clean          Remove generated files"

# ============================
# Mock Generation
# ============================
mock:
	@echo ">> Generating mocks..."
	@mkdir -p $(MOCK_DIR)
	$(MOCKGEN) \
		-source=$(EXECUTOR_SRC) \
		-destination=$(MOCK_FILE) \
		-package=mocks

# ============================
# Unit Tests
# ============================
test:
	@echo ">> Running repository tests..."
	$(GO) test -v $(REPO_PKG)

test-all:
	@echo ">> Running all tests..."
	$(GO) test -v ./...

# ============================
# Coverage
# ============================
coverage:
	@echo ">> Running tests with coverage..."
	$(GO) test ./... -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-html:
	@echo ">> Opening coverage report..."
	$(GO) tool cover -html=$(COVERAGE_FILE)

# ============================
# Clean
# ============================
clean:
	@echo ">> Cleaning generated files..."
	@rm -f $(MOCK_FILE)
	@rm -f $(COVERAGE_FILE)
