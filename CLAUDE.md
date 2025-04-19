# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Build and run: `go run cmd/app/main.go`
- Run all tests: `go test -race -cover ./...`
- Run specific test: `go test -race -cover ./path/to/package -run TestName`
- Linting: `golangci-lint run ./...` or `make lint`
- Formatting: `golangci-lint fmt ./...` or `make fmt`

## Code Style Guidelines
- **Structure**: Follow Domain-Driven Design with clean architecture layers
- **Imports**: Standard lib → Third-party → Project (enforced by goimports)
- **Formatting**: 2-space indentation, LF line endings, trim trailing whitespace
- **Naming**: PascalCase for exported, camelCase for non-exported identifiers
- **Error Handling**: Use domain/domainerror package for typed errors, always wrap with context
- **Testing**: Table-driven tests with testify/assert, mocks with testify/mock
- **Types**: Use constructor functions (NewXxx) for complex types
- **File Organization**: entity.go, repository.go pattern in domain packages

## Database Operations
- Start database: `docker-compose up`
- Initialize schema: Run scripts in docker/mysql/initdb.d/