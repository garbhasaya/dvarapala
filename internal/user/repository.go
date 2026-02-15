package user

import (
	"context"
	"fmt"
	"log/slog"

	"dvarapala/ent"
	"dvarapala/ent/user"
)

// UserRepository handles database operations for users.
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	created, err := r.client.User.
		Create().
		SetAppID(u.AppID).
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
	return r.GetByID(ctx, created.ID)
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := r.client.User.Query().
		Where(user.IDEQ(id)).
		WithApp().
		Only(ctx)
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
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := r.client.User.Query().
		Where(user.EmailEQ(email)).
		WithApp().
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
func (r *UserRepository) List(ctx context.Context) ([]*ent.User, error) {
	users, err := r.client.User.Query().WithApp().All(ctx)
	if err != nil {
		slog.Error("database error: failed to list users", "error", err)
		return nil, err
	}
	return users, nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, id int, u *ent.User) (*ent.User, error) {
	updated, err := r.client.User.UpdateOneID(id).
		SetAppID(u.AppID).
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
	return r.GetByID(ctx, updated.ID)
}

// Delete deletes a user by their ID.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
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
