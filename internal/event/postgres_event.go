package event

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/db/query"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db,
	}
}

// List returns all events
func (r *Repository) List(ctx context.Context) ([]*persistence.Event, error) {
	q := query.Use(r.db)
	rows, err := q.WithContext(ctx).Event.Find()
	if err != nil {
		return nil, fmt.Errorf("repo: list events %w", err)
	}

	return rows, nil
}

// ListByUserID returns all events by userID
func (r *Repository) ListByUserID(ctx context.Context) ([]*persistence.Event, error) {
	userID, _ := auth.UserIDFromContext(ctx)
	q := query.Use(r.db)
	rows, err := q.WithContext(ctx).Event.Where(q.Event.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, fmt.Errorf("repo: list events %w", err)
	}

	return rows, nil
}

// Create creates a new event
func (r *Repository) Create(ctx context.Context, event persistence.Event) (*persistence.Event, error) {
	q := query.Use(r.db)
	err := q.WithContext(ctx).Event.Create(&event)
	if err != nil {
		return nil, fmt.Errorf("repo: create event %w", err)
	}

	return &event, nil
}

// Get returns an event by id
func (r *Repository) Get(ctx context.Context, id string) (*persistence.Event, error) {
	q := query.Use(r.db)
	event, err := q.WithContext(ctx).Event.Where(q.Event.ID.Eq(id)).First()
	if err != nil {
		return nil, fmt.Errorf("repo: get event %w", err)
	}

	return event, nil
}

// Update updates an event
func (r *Repository) Update(ctx context.Context, event persistence.Event) error {
	q := query.Use(r.db)
	return q.WithContext(ctx).Event.Save(&event)
}

// Delete deletes an event
func (r *Repository) Delete(ctx context.Context, id string) error {
	q := query.Use(r.db)
	_, err := q.WithContext(ctx).Event.Where(q.Event.ID.Eq(id)).Delete()
	return err
}
