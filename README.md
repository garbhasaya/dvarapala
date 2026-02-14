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
- `make clean`: Deep clean of containers, images, and volumes.

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

## API endpoints

- POST /users/auth
- POST /users/password/forget
- POST /users/password/reset
- POST /users
- GET /users
- POST /users/{id}
- DELETE /users/{id}
- GET /users/{id}
