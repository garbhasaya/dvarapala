.PHONY: build up down restart logs ps test lint swag clean shell

# Docker Compose commands
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

ps:
	docker-compose ps

# Run tests inside the container
test:
	docker-compose run --rm api go test -v ./...

# Run linter using a docker container
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

# Generate Swagger documentation
# Note: This assumes swag is installed in the container or runs via a separate image
swag:
	docker-compose run --rm api go install github.com/swaggo/swag/cmd/swag@latest
	docker-compose run --rm api swag init -g cmd/api/main.go

# Open a shell in the running api container
shell:
	docker-compose exec api sh

# Database migrations
migrate-gen:
	docker run --rm -v $(shell pwd):/app -w /app \
		-e CGO_ENABLED=1 \
		-e CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
		golang:1.26-alpine \
		sh -c "apk add --no-cache build-base && go run -mod=mod ent/migrate/main.go $(name)"

# Clean up containers, images, and volumes
clean:
	docker-compose down --rmi all --volumes --remove-orphans
