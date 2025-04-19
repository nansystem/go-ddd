.PHONY: lint format test

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

test:
	go test -race -cover ./...

install-hooks:
	cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

migrate-up:
	docker exec -it go-ddd-mysql mysql -u ddduser -pdddpass -e "CREATE DATABASE IF NOT EXISTS go_ddd;"
	docker exec -it go-ddd-mysql mysql -u ddduser -pdddpass -e "USE go_ddd; source /docker-entrypoint-initdb.d/01_schema.sql;"
	docker exec -it go-ddd-mysql mysql -u ddduser -pdddpass -e "USE go_ddd; source /docker-entrypoint-initdb.d/02_testdata.sql;"

migrate-down:
	docker exec -it go-ddd-mysql mysql -u ddduser -pdddpass -e "DROP DATABASE IF EXISTS go_ddd;"
