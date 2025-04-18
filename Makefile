.PHONY: lint format test

lint: ## コードの静的解析を実行します
	golangci-lint run ./...

format: ## コードのフォーマットを整形します
	golangci-lint fmt ./...

test:
	go test -race -cover ./...

install-hooks:
	cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
