package user

import (
	"context"
	"testing"
	"time"

	"dvarapala/ent/enttest"
	"dvarapala/internal/platform/auth"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() {
		err := client.Close()
		assert.NoError(t, err)
	}()

	repo := NewUserRepository(client)
	jwtManager := auth.NewJWTManager("secret", time.Hour)
	svc := NewUserService(repo, jwtManager)

	ctx := context.Background()
	req := CreateUserRequest{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	u, err := svc.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, req.Email, u.Email)
	assert.Equal(t, req.Firstname, u.Firstname)
	assert.Equal(t, req.Lastname, u.Lastname)
	assert.Equal(t, int8(1), u.Status)
}

func TestService_Authenticate(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent_auth?mode=memory&cache=shared&_fk=1")
	defer func() {
		err := client.Close()
		assert.NoError(t, err)
	}()

	repo := NewUserRepository(client)
	jwtManager := auth.NewJWTManager("secret", time.Hour)
	svc := NewUserService(repo, jwtManager)

	ctx := context.Background()
	email := "auth@example.com"
	password := "password123"

	_, err := svc.Create(ctx, CreateUserRequest{
		Firstname: "Auth",
		Lastname:  "User",
		Email:     email,
		Password:  password,
	})
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		res, err := svc.Authenticate(ctx, AuthRequest{
			Email:    email,
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Token)
		assert.Equal(t, email, res.User.Email)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		res, err := svc.Authenticate(ctx, AuthRequest{
			Email:    email,
			Password: "wrongpassword",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		res, err := svc.Authenticate(ctx, AuthRequest{
			Email:    "nonexistent@example.com",
			Password: password,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

func TestService_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent_update?mode=memory&cache=shared&_fk=1")
	defer func() {
		err := client.Close()
		assert.NoError(t, err)
	}()

	repo := NewUserRepository(client)
	jwtManager := auth.NewJWTManager("secret", time.Hour)
	svc := NewUserService(repo, jwtManager)

	ctx := context.Background()
	u, err := svc.Create(ctx, CreateUserRequest{
		Firstname: "Original",
		Lastname:  "Name",
		Email:     "original@example.com",
		Password:  "password123",
	})
	assert.NoError(t, err)

	newName := "Updated"
	newStatus := int8(0)
	req := UpdateUserRequest{
		Firstname: &newName,
		Status:    &newStatus,
	}

	updated, err := svc.Update(ctx, u.ID, req)
	assert.NoError(t, err)
	assert.Equal(t, newName, updated.Firstname)
	assert.Equal(t, newStatus, updated.Status)
	assert.Equal(t, u.Email, updated.Email) // Should remain unchanged
}

func TestService_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent_delete?mode=memory&cache=shared&_fk=1")
	defer func() {
		err := client.Close()
		assert.NoError(t, err)
	}()

	repo := NewUserRepository(client)
	jwtManager := auth.NewJWTManager("secret", time.Hour)
	svc := NewUserService(repo, jwtManager)

	ctx := context.Background()
	u, err := svc.Create(ctx, CreateUserRequest{
		Firstname: "Delete",
		Lastname:  "Me",
		Email:     "delete@example.com",
		Password:  "password123",
	})
	assert.NoError(t, err)

	err = svc.Delete(ctx, u.ID)
	assert.NoError(t, err)

	// Verify it's gone
	_, err = svc.GetByID(ctx, u.ID)
	assert.Error(t, err)
}
