package app

import (
	"time"
)

// App represents the domain model for an app.
type App struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateAppRequest defines the payload for creating an app.
type CreateAppRequest struct {
	Name   string `json:"name" validate:"required"`
	Status int8   `json:"status" validate:"omitempty"`
}

// UpdateAppRequest defines the payload for updating an app.
type UpdateAppRequest struct {
	Name   *string `json:"name" validate:"omitempty"`
	Status *int8   `json:"status" validate:"omitempty"`
}
