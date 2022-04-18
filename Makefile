create-bank-network:
	docker network create banking-system-network

setup-postgres:
	docker run --name postgres-local --network banking-system-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:12.10-alpine 

create-db:
		docker exec -it postgres-local createdb --username=postgres --owner=postgres bank
		docker exec -it postgres-local createdb --username=postgres --owner=postgres bank_test

create-migration:
	migrate create -ext sql -dir db/migration -seq -digits 1 $(migration_name)

download:
	go mod download

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

mock-db:
	mockgen -package mockdb -destination db/mock/store.go github.com/skamranahmed/banking-system/db/sqlc Store

build:
	docker build -t banking-system:latest .

dockerized-server-run:
	docker run --name banking-system --network banking-system-network -p 8080:8080 -e DB_HOST="postgresql://postgres:password@postgres-local:5432/bank?sslmode=disable" banking-system:latest

dc-up:
	docker-compose up -d

dc-down:
	docker-compose down

run:
	go run main.go

.PHONY: create-migration migrate-up migrate-up-test migrate-down migrate-down-test sqlc-gen test mock-db build run	