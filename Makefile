.PHONY: lint format test

lint:
	golangci-lint run ./...

format:
	go fmt ./...
	go tool goimports -local github.com/nansystem/go-ddd -w .

test:
	go test -race -cover ./...

install-hooks:
	cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
