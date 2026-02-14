package user

import (
	"context"
	"fmt"
	"log/slog"

	"dvarapala/ent"
	"dvarapala/ent/user"
)

// Repository handles database operations for users.
type Repository struct {
	client *ent.Client
}

// NewRepository creates a new user repository.
func NewRepository(client *ent.Client) *Repository {
	return &Repository{client: client}
}

// Create creates a new user in the database.
func (r *Repository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	created, err := r.client.User.
		Create().
		SetFirstname(u.Firstname).
		SetLastname(u.Lastname).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		Save(ctx)
	if err != nil {
		slog.Error("database error: failed to create user", "email", u.Email, "error", err)
		return nil, err
	}
	return created, nil
}

// GetByID retrieves a user by their ID.
func (r *Repository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("user not found in database", "id", id)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		slog.Error("database error: failed to get user by id", "id", id, "error", err)
		return nil, err
	}
	return u, nil
}

// GetByEmail retrieves a user by their email.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("user not found in database", "email", email)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		slog.Error("database error: failed to get user by email", "email", email, "error", err)
		return nil, err
	}
	return u, nil
}

// List retrieves all users.
func (r *Repository) List(ctx context.Context) ([]*ent.User, error) {
	users, err := r.client.User.Query().All(ctx)
	if err != nil {
		slog.Error("database error: failed to list users", "error", err)
		return nil, err
	}
	return users, nil
}

// Update updates an existing user.
func (r *Repository) Update(ctx context.Context, id int, u *ent.User) (*ent.User, error) {
	updated, err := r.client.User.UpdateOneID(id).
		SetFirstname(u.Firstname).
		SetLastname(u.Lastname).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		Save(ctx)
	if err != nil {
		slog.Error("database error: failed to update user", "id", id, "error", err)
		return nil, err
	}
	return updated, nil
}

// Delete deletes a user by their ID.
func (r *Repository) Delete(ctx context.Context, id int) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.Warn("user not found for deletion", "id", id)
			return fmt.Errorf("user not found: %w", err)
		}
		slog.Error("database error: failed to delete user", "id", id, "error", err)
		return err
	}
	return nil
}
