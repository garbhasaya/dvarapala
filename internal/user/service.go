package user

import (
	"context"
	"errors"
	"fmt"

	"dvarapala/ent"
	"dvarapala/internal/platform/auth"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the business logic for users.
type Service interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, id int, req UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id int) error
	Authenticate(ctx context.Context, req AuthRequest) (*AuthResponse, error)
}

type service struct {
	repo *Repository
	jwt  *auth.JWTManager
}

// NewService creates a new user service.
func NewService(repo *Repository, jwt *auth.JWTManager) Service {
	return &service{repo: repo, jwt: jwt}
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &ent.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Status:    1,
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("repository create: %w", err)
	}

	return s.toDomain(created), nil
}

func (s *service) GetByID(ctx context.Context, id int) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toDomain(u), nil
}

func (s *service) List(ctx context.Context) ([]*User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]*User, len(users))
	for i, u := range users {
		domainUsers[i] = s.toDomain(u)
	}
	return domainUsers, nil
}

func (s *service) Update(ctx context.Context, id int, req UpdateUserRequest) (*User, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Firstname != nil {
		existing.Firstname = *req.Firstname
	}
	if req.Lastname != nil {
		existing.Lastname = *req.Lastname
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		existing.Password = string(hashedPassword)
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	updated, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	return s.toDomain(updated), nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) Authenticate(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(u.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  *s.toDomain(u),
	}, nil
}

func (s *service) toDomain(u *ent.User) *User {
	return &User{
		ID:        u.ID,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
