package app

import (
	"context"
	"fmt"
	"log/slog"

	"dvarapala/ent"
)

// AppRepository handles database operations for apps.
type AppRepository struct {
	client *ent.Client
}

// NewAppRepository creates a new app repository.
func NewAppRepository(client *ent.Client) *AppRepository {
	return &AppRepository{client: client}
}

// Create creates a new app in the database.
func (r *AppRepository) Create(ctx context.Context, a *ent.App) (*ent.App, error) {
	created, err := r.client.App.
		Create().
		SetName(a.Name).
		SetStatus(a.Status).
		Save(ctx)
	if err != nil {
		slog.Error("database error: failed to create app", "name", a.Name, "error", err)
		return nil, err
	}
	return created, nil
}

// GetByID retrieves an app by its ID.
func (r *AppRepository) GetByID(ctx context.Context, id int) (*ent.App, error) {
	a, err := r.client.App.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("app not found in database", "id", id)
			return nil, fmt.Errorf("app not found: %w", err)
		}
		slog.Error("database error: failed to get app by id", "id", id, "error", err)
		return nil, err
	}
	return a, nil
}

// List retrieves all apps.
func (r *AppRepository) List(ctx context.Context) ([]*ent.App, error) {
	apps, err := r.client.App.Query().All(ctx)
	if err != nil {
		slog.Error("database error: failed to list apps", "error", err)
		return nil, err
	}
	return apps, nil
}

// Update updates an existing app.
func (r *AppRepository) Update(ctx context.Context, id int, a *ent.App) (*ent.App, error) {
	updated, err := r.client.App.UpdateOneID(id).
		SetName(a.Name).
		SetStatus(a.Status).
		Save(ctx)
	if err != nil {
		slog.Error("database error: failed to update app", "id", id, "error", err)
		return nil, err
	}
	return updated, nil
}

// Delete deletes an app by its ID.
func (r *AppRepository) Delete(ctx context.Context, id int) error {
	err := r.client.App.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("app not found for deletion", "id", id)
			return fmt.Errorf("app not found: %w", err)
		}
		slog.Error("database error: failed to delete app", "id", id, "error", err)
		return err
	}
	return nil
}
