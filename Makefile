include .env
export 

.PHONY: clean codegen test mocks dev dev-db dev-db-down reset migrate migration test-db test-db-down migrate-test reset-test

clean:
	rm main

codegen:
	sqlc generate

test:
	go test ./... 

mocks:
	mockgen -package=mocks -source=internal/service/interface.go -destination=internal/test/mocks/service.go

dev:
	go build ./cmd/go-api-template/main.go && ./main

dev-db:
	docker compose -p go-api-template -f docker-compose.yml up --detach

dev-db-down:
	docker compose -p go-api-template -f docker-compose.yml down -v

reset:
	make dev-db-down && make dev-db

migrate:
	migrate -database "$(DATABASE_URL)?sslmode=disable" -path ./db/migrations up

# Usage: make migration name=your_migration_name
migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

test-db:
	docker compose -p go-api-template-test -f docker-compose.yml up --detach

test-db-down:
	docker compose -p go-api-template-test -f docker-compose.test.yml down -v

migrate-test:
	migrate -database "$(TEST_DATABASE_URL)?sslmode=disable" -path ./db/migrations up

reset-test:
	make test-db-down && make test-db && make migrate-test