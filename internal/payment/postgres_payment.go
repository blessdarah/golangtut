package payment

import (
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
		db: db,
	}
}

func (r *Repository) List(
	ctx context.Context,
	filters map[string]string,
) ([]*persistence.Payment, error) {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0, len(filters))
	for col, val := range filters {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", col), val)
		})
	}

	r.db.Logger.Info(ctx, "repo payment scopes", scopes)

	var payments []*persistence.Payment

	err := r.db.WithContext(ctx).
		Model(&persistence.Payment{}).
		Scopes(scopes...).
		Find(&payments).Error

	if err != nil {
		return nil, fmt.Errorf("repo: list payments %w", err)
	}

	return payments, nil
}

func (r *Repository) Create(
	ctx context.Context,
	payment *persistence.Payment,
) (*persistence.Payment, error) {
	q := query.Use(r.db)

	err := q.WithContext(ctx).Payment.Create(payment)
	if err != nil {
		return nil, fmt.Errorf("repo: create payment %w", err)
	}

	return payment, nil
}
