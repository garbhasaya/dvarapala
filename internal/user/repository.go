package user

import (
	"context"
	"fmt"

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
	return r.client.User.
		Create().
		SetFirstname(u.Firstname).
		SetLastname(u.Lastname).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		Save(ctx)
}

// GetByID retrieves a user by their ID.
func (r *Repository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
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
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}
	return u, nil
}

// List retrieves all users.
func (r *Repository) List(ctx context.Context) ([]*ent.User, error) {
	return r.client.User.Query().All(ctx)
}

// Update updates an existing user.
func (r *Repository) Update(ctx context.Context, id int, u *ent.User) (*ent.User, error) {
	return r.client.User.UpdateOneID(id).
		SetFirstname(u.Firstname).
		SetLastname(u.Lastname).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		Save(ctx)
}

// Delete deletes a user by their ID.
func (r *Repository) Delete(ctx context.Context, id int) error {
	err := r.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user not found: %w", err)
		}
		return err
	}
	return nil
}
