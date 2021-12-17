postgres:
	docker run --name postgres14-eshop -p 5555:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14-eshop createdb --username=postgres --owner=postgres eshop

dropdb:
	docker exec -it postgres14-eshop dropdb eshop

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@192.168.10.139:5555/eshop?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:secret@192.168.10.139:5555/eshop?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@192.168.10.139:5555/eshop?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:secret@192.168.10.139:5555/eshop?sslmode=disable" -verbose down 1

cimigrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5555/eshop?sslmode=disable" -verbose up

cimigratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5555/eshop?sslmode=disable" -verbose down

sqlc:
	sqlc generate

sqlcwin:
	docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate

sqlcfix:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mnbenghuzzi/simplebank/v2/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 cimigrateup cimigratedown sqlc sqlcfix sqlcwin test server mock