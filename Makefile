# Build the application
build:
	@go build -o bin/RealState cmd/main.go

# Run tests
test:
	@go test -v ./...

# Run the application
run: build
	@./bin/RealState

# Create a new migration
migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))


# Apply all migrations (up)
migrate-up:
	@go run cmd/migrate/main.go up

# Rollback the last migration (down)
migrate-down:
	@go run cmd/migrate/main.go down
