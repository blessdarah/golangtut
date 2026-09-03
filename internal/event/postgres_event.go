package event

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/model"
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
func (r *Repository) List(ctx context.Context) ([]*model.Event, error) {
	var events []*model.Event
	err := r.db.WithContext(ctx).Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("repo: list events %w", err)
	}

	return events, nil
}

// ListByUserID returns all events by userID
func (r *Repository) ListByUserID(ctx context.Context) ([]*model.Event, error) {
	var events []*model.Event
	userID, _ := auth.UserIDFromContext(ctx)

	err := r.db.WithContext(ctx).
		Where("user_id = ? ", userID).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("repo: list events %w", err)
	}

	return events, nil
}

// Create creates a new event
func (r *Repository) Create(ctx context.Context, event model.Event) (*model.Event, error) {
	err := r.db.WithContext(ctx).Create(&event).Error
	if err != nil {
		return nil, fmt.Errorf("repo: create event %w", err)
	}

	return &event, nil
}

// Get returns an event by id
func (r *Repository) Get(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).
		Where("id = ? ", id).
		First(&event).Error
	if err != nil {
		return nil, fmt.Errorf("repo: get event %w", err)
	}

	return &event, nil
}

// Update updates an event
func (r *Repository) Update(ctx context.Context, event model.Event) error {
	return r.db.WithContext(ctx).Save(&event).Error
}

// Delete deletes an event
func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ? ", id).
		Delete(&model.Event{}).Error
}
