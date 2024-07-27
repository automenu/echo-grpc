.PHONY: help
help:
	@echo "---"
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  dev         Start dev server"
	@echo "  ping        Ping the server using gRPC client"
	@echo "  check       Run all checks"
	@echo "  lintfix     Run linters and fix some issues"
	@echo "  codegen     Generate code"
	@echo "---"

.PHONY: dev
dev:
	@echo "Starting dev server..."
	go run ./cmd/server/main.go

.PHONY: ping
ping:
	@echo "Pinging the server using gRPC client..."
	go run ./cmd/client/main.go

.PHONY: check
check: lint check-codegen check-breaking check-tidy
	@echo "Checks done"

.PHONY: lint
lint:
	@echo "Linting..."
	golangci-lint run
	buf lint
	buf format -d --exit-code
	npx --yes prettier --check .
	@echo "Linting done"

.PHONY: lintfix
lintfix:
	@echo "Linting and fixing some linting issues..."
	golangci-lint run --fix
	buf format -w
	npx --yes prettier --write .
	@echo "Linting and fixing done"

.PHONY: codegen
codegen:
	@echo "Generating code..."
	buf generate
	@echo "Code generation done"

.PHONY: check-codegen
check-codegen:
	@echo "Checking codegen..."
	buf generate
	test -z "$$(git status --porcelain | tee /dev/stderr)"
	@echo "Codegen is up to date"

.PHONY: check-breaking
check-breaking:
	@echo "Checking if codegen is breaking..."
	buf breaking --against 'https://github.com/automenu/echo-grpc.git#branch=main'
	@echo "Codegen is not breaking"

.PHONY: check-tidy
check-tidy:
	@echo "Checking if 'go mod tidy' is needed..."
	go mod tidy
	test -z "$$(git status --porcelain | tee /dev/stderr)"
	@echo "Check tidy done"
