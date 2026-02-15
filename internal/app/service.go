package app

import (
	"context"
	"fmt"
	"log/slog"

	"dvarapala/ent"
)

// AppService defines the business logic for apps.
type AppService interface {
	Create(ctx context.Context, req CreateAppRequest) (*App, error)
	GetByID(ctx context.Context, id int) (*App, error)
	List(ctx context.Context) ([]*App, error)
	Update(ctx context.Context, id int, req UpdateAppRequest) (*App, error)
	Delete(ctx context.Context, id int) error
}

type appService struct {
	repo *AppRepository
}

// NewAppService creates a new app service.
func NewAppService(repo *AppRepository) AppService {
	return &appService{repo: repo}
}

func (s *appService) Create(ctx context.Context, req CreateAppRequest) (*App, error) {
	slog.Info("creating app", "name", req.Name)

	status := int8(1)
	if req.Status != 0 {
		status = req.Status
	}

	a := &ent.App{
		Name:   req.Name,
		Status: status,
	}

	created, err := s.repo.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("repository create: %w", err)
	}

	slog.Info("app created successfully", "id", created.ID, "name", created.Name)
	return s.toDomain(created), nil
}

func (s *appService) GetByID(ctx context.Context, id int) (*App, error) {
	slog.Info("getting app by id", "id", id)
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toDomain(a), nil
}

func (s *appService) List(ctx context.Context) ([]*App, error) {
	slog.Info("listing apps")
	apps, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	domainApps := make([]*App, len(apps))
	for i, a := range apps {
		domainApps[i] = s.toDomain(a)
	}
	return domainApps, nil
}

func (s *appService) Update(ctx context.Context, id int, req UpdateAppRequest) (*App, error) {
	slog.Info("updating app", "id", id)
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	updated, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	slog.Info("app updated successfully", "id", id)
	return s.toDomain(updated), nil
}

func (s *appService) Delete(ctx context.Context, id int) error {
	slog.Info("deleting app", "id", id)
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	slog.Info("app deleted successfully", "id", id)
	return nil
}

func (s *appService) toDomain(a *ent.App) *App {
	return &App{
		ID:        a.ID,
		Name:      a.Name,
		Status:    a.Status,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
