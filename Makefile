AWS_DB_URL=postgresql://root:simplebanksecret@simple-bank.c6pqv8upf7jt.us-east-1.rds.amazonaws.com:5433/simple_bank
LOCAL_DB_URL=postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(AWS_DB_URL)" -verbose up 

migrateup1:
	migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(AWS_DB_URL)" -verbose down 

migratedown1:
	migrate -path db/migration -database "$(LOCAL_DB_URL)" -verbose down 1

sqlc:
	docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go
	
mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/fredrick/simplebank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 db_docs db_schema