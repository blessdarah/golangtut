package event

import (
	"blessdarah/tuts/internal/model"
	"context"
)

type repository interface {
	List(ctx context.Context) ([]*model.Event, error)
	ListByUserID(ctx context.Context) ([]*model.Event, error)
	Create(ctx context.Context, event model.Event) (*model.Event, error)
	Get(ctx context.Context, id string) (*model.Event, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo,
	}
}

// GetAll returns all events
func (s *Service) GetAll(ctx context.Context) ([]*model.Event, error) {
	return s.repo.List(ctx)
}

// GetByUserID returns all events by userID
func (s *Service) GetByUserID(ctx context.Context) ([]*model.Event, error) {
	return s.repo.ListByUserID(ctx)
}

// Create creates a new event
func (s *Service) Create(ctx context.Context, event model.Event) (*model.Event, error) {
	return s.repo.Create(ctx, event)
}

// Get returns an event by id
func (s *Service) Get(ctx context.Context, id string) (*model.Event, error) {
	return s.repo.Get(ctx, id)
}

// Update updates an event
func (s *Service) Update(ctx context.Context, event model.Event) error {
	return s.repo.Update(ctx, event)
}

// Delete deletes an event
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
