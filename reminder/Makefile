run:
	go run ./cmd/service/

tidy:
	go mod tidy

MIGRATIONS_PATH=db/migrations
DB_URL=sqlite3://reminders.db
.PHONY: migrate-up migrate-down migrate-create migrate-version
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down $(or $(N),1)
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version
