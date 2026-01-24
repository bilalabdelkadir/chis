
# Run the server
run-api:
	go run cmd/api/main.go
# Run the worker
run-worker:
	go run cmd/worker/main.go

run-grpc:
	go run cmd/delivery/main.go
# Database migrations
migrate-up:
	migrate -path internal/database/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path internal/database/migrations -database "$(DB_URL)" down

migrate-version:
	migrate -path internal/database/migrations -database "$(DB_URL)" version

# Create a new migration (usage: make migrate-create name=add_something)
migrate-create:
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/delivery/delivery.proto