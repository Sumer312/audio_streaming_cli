include .env

run:
	@go run main.go

migrateUp:
	@goose -dir ./sql/schems $(DB_TYPE) $(DB_URL) up

migrateDown:
	@goose -dir ./sql/schemas $(DB_TYPE) $(DB_URL) down
