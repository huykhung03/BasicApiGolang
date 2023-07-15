postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=oneanhiuemlove33 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_shop

dropdb:
	docker exec -it postgres12 dropdb simple_shop

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:oneanhiuemlove33@localhost:5432/simple_shop?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:oneanhiuemlove33@localhost:5432/simple_shop?sslmode=disable" -verbose down
	
sqlc:
	docker run --rm -v "C:\SimpleShop:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: 
	postgres 
	createdb dropdb 
	migrateup migratedown
	sqlc
