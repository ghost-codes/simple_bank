postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:15

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=root bank

dropdb:
	docker exec -it postgres15 dropdb bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: migrateup migratedown sqlc test server