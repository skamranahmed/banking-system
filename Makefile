create-migration:
	migrate create -ext sql -dir db/migration -seq -digits 1 <migration_name>

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

.PHONY: create-migration migrate-up