migrate_up:
	goose -dir="./database/schema" sqlite3 ./chat.db up
migrate_down:
	goose -dir="./database/schema" sqlite3 ./chat.db down
init_db:
	 goose -dir="./database/schema" sqlite3 ./chat.db create init sql
drop_db:
	rm ./chat.db
generate:
	sqlc generate
test:
	go test -v -cover ./...
