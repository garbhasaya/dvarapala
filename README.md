# Project overview

# Architecture

## Go packages

- `ent` (https://github.com/ent/ent) ORM
- `chi` (https://github.com/go-chi/chi) for routing
- `testify` (https://github.com/stretchr/testify) for writing and running unit tests
- `viper` (https://github.com/spf13/viper) manage multiple environments i.e dev, test, CAT, prod configurations
- `slog` go standard library for logging
- `swag` (https://github.com/swaggo/swag) generate RESTful API documentation.
- `golangci-lint` (https://github.com/golangci/golangci-lint) linter
- `jwt` (https://github.com/golang-jwt/jwt) implementation of JSON Web Tokens (JWT)
- `cors` (https://github.com/go-chi/cors) CORS net/http middleware for Go
- `httprate` (https://github.com/go-chi/httprate) net/http rate limiter middleware
- `validator` (https://github.com/go-playground/validator) field validation, including Cross Field, Cross Struct, Map, Slice and Array diving

## Directory structure

```
/cmd/api/main.go
/internal/
  user/
    handler.go
    service.go
    repository.go
    model.go
  platform/
    http/
    middleware.go
    response.go
    router.go
  db/
    sqlite.go
/pkg/
  logger/
  config/
```

## Code architecture

### Use directional dependencies:
HTTP → Service → Repository

#### Handler (Delivery Layer):
- Only HTTP concerns
- No business logic
```
type UserHandler struct {
    svc UserService
}
```

#### Service (Business Logic):
- Pure Go logic
- No HTTP, no SQL
```
type UserService interface {
    Create(ctx context.Context, u User) error
}
```

#### Repository (Persistence)
- DB logic only
- Implements interfaces


## Design patterns

- Dependancy injection (DI)
- Interface Segregation (Very Important in Go)
```
type UserWriter interface {
    Save(ctx context.Context, u User) error
}

type UserReader interface {
    FindByID(ctx context.Context, id string) (User, error)
}
```
- Error Handling Pattern (No Exceptions)
Sentinel + Wrapped Errors
```
var ErrUserNotFound = errors.New("user not found")

if err != nil {
    return fmt.Errorf("create user: %w", err)
}
```
Translate errors at the boundary (HTTP)
```
if errors.Is(err, ErrUserNotFound) {
    http.Error(w, "not found", http.StatusNotFound)
}
```
- Context Propagation (Mandatory)
```
func (s *service) Create(ctx context.Context, u User) error
```

## Requirement

- Go v1.26
- SQLite v3.51.2

## Development

The project uses Docker and a Makefile for development.

- `make build`: Build the Docker images.
- `make up`: Start the containers in the background.
- `make down`: Stop and remove the containers.
- `make restart`: Restart the services.
- `make logs`: Follow the container logs.
- `make ps`: List the running containers.
- `make test`: Run all Go tests inside the container.
- `make lint`: Run `golangci-lint` using a dedicated Docker image.
- `make swag`: Generate Swagger documentation.
- `make shell`: Open an interactive shell inside the API container.
- `make migrate-gen name=migration_name`: Generate a new versioned migration file.
- `make migrate-apply`: Apply all pending migrations to the database.
- `make clean`: Deep clean of containers, images, and volumes.

## Database Migrations

This project uses **Ent** with **Atlas** for versioned migrations. Follow these steps when you need to change the database schema:

### 1. Create or Modify the Schema

#### To create a new table:
Initialize a new schema file:
```bash
docker run --rm -v $(pwd):/app -w /app golang:1.26-alpine go run -mod=mod entgo.io/ent/cmd/ent new TableName
```
Then define the fields in `ent/schema/tablename.go`.

#### To modify an existing table:
Update the schema definitions in the `ent/schema/` directory (e.g., `ent/schema/user.go`).

### 2. Generate Ent Code
After modifying the schema, regenerate the Ent runtime code:
```bash
docker run --rm -v $(pwd):/app -w /app golang:1.26-alpine go generate ./ent/...
```

### 3. Generate Migration Files
Generate a new SQL migration file by comparing your schema changes against an in-memory database:
```bash
make migrate-gen name=add_new_field_to_user
```
This will create new `.sql` files in `ent/migrate/migrations/`.

### 4. Apply Migrations
You can manually apply migrations to the database using:
```bash
make migrate-apply
```

Additionally, in the current development setup, the application automatically applies migrations on startup using `client.Schema.Create` in `internal/db/sqlite.go`. You can restart the service to trigger this:
```bash
make restart
```

## Database Persistence

The SQLite database is stored at `/app/data/dvarapala.db` inside the container. This path is persisted using a bind mount to the local `./data` directory in the project root.

- **Host Path**: `./data/dvarapala.db`
- **Container Path**: `/app/data/dvarapala.db`
- **Environment Variable**: `DB_PATH`

The database initialization is fully aligned with the Ent migration setup. On every startup, the application verifies the schema against the generated Ent code and applies any necessary changes to the SQLite file, ensuring the physical database always matches your versioned migration files.

## Database schema

### user

- ID - int - primary key - auto increment
- Firstname
- Lastname
- Email
- Password
- Status - smallint - 0 or 1
- Created at
- Updated at

## Service URLs

- **API Gateway**: [http://localhost:8080](http://localhost:8080)
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API endpoints

- POST /users/auth
- POST /users/password/forget
- POST /users/password/reset
- POST /users
- GET /users
- POST /users/{id}
- DELETE /users/{id}
- GET /users/{id}
