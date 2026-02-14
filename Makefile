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
	docker run --rm -v $(shell pwd):/app -w /app \
		-e CGO_ENABLED=1 \
		-e CGO_CFLAGS="-D_LARGEFILE64_SOURCE" \
		golang:1.26-alpine \
		sh -c "apk add --no-cache build-base && go test -v ./..."

# Run linter using a docker container
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

# Generate Swagger documentation
swag:
	docker run --rm -v $(shell pwd):/app -w /app golang:latest sh -c "go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go --parseDependency --parseInternal"

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

migrate-apply:
	docker-compose run --rm atlas migrate apply \
		--url "sqlite:///data/dvarapala.db?_fk=1" \
		--dir "file://ent/migrate/migrations" \
		--allow-dirty

# Clean up containers, images, and volumes
clean:
	docker-compose down --rmi all --volumes --remove-orphans
