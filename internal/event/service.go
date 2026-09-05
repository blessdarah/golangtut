package event

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	List(ctx context.Context) ([]*persistence.Event, error)
	ListByUserID(ctx context.Context) ([]*persistence.Event, error)
	Create(ctx context.Context, event persistence.Event) (*persistence.Event, error)
	Get(ctx context.Context, id string) (*persistence.Event, error)
	Update(ctx context.Context, event persistence.Event) error
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
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]*model.Event, len(rows))
	for i, row := range rows {
		e := toDomainEvent(row)
		events[i] = &e
	}

	return events, nil
}

// GetByUserID returns all events by userID
func (s *Service) GetByUserID(ctx context.Context) ([]*model.Event, error) {
	rows, err := s.repo.ListByUserID(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]*model.Event, len(rows))
	for i, row := range rows {
		e := toDomainEvent(row)
		events[i] = &e
	}

	return events, nil
}

// Create creates a new event
func (s *Service) Create(ctx context.Context, event model.Event) (*model.Event, error) {
	row := toSchemaEvent(event)
	row.ID = uuid.NewString()

	created, err := s.repo.Create(ctx, *row)
	if err != nil {
		return nil, err
	}

	res := toDomainEvent(created)
	return &res, nil
}

// Get returns an event by id
func (s *Service) Get(ctx context.Context, id string) (*model.Event, error) {
	row, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	res := toDomainEvent(row)
	return &res, nil
}

// Update updates an event
func (s *Service) Update(ctx context.Context, event model.Event) error {
	row := toSchemaEvent(event)
	return s.repo.Update(ctx, *row)
}

// Delete deletes an event
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toDomainEvent(e *persistence.Event) model.Event {
	userID := e.UserID
	return model.Event{
		ID:          e.ID,
		UserID:      &userID,
		Name:        e.Name,
		Description: e.Description,
		Venue:       e.Venue,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toSchemaEvent(e model.Event) *persistence.Event {
	userID := ""
	if e.UserID != nil {
		userID = *e.UserID
	}

	return &persistence.Event{
		ID:          e.ID,
		UserID:      userID,
		Name:        e.Name,
		Description: e.Description,
		Venue:       e.Venue,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
