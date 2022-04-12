create-migration:
	migrate create -ext sql -dir db/migration -seq -digits 1 <migration_name>

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate-up-test:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank_test?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

migrate-down-test:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank_test?sslmode=disable" -verbose down

sqlc-gen:
	sqlc generate

test:
	go clean -testcache && go test -v -cover ./...

.PHONY: create-migration migrate-up migrate-up-test migrate-down migrate-down-test sqlc-gen test