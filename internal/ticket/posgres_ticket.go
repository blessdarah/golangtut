package ticket

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

func (r *Repository) List(ctx context.Context) ([]*persistence.Ticket, error) {
	q := query.Use(r.db)
	rows, err := q.WithContext(ctx).Ticket.Find()
	if err != nil {
		return nil, fmt.Errorf("repo: list tickets %w", err)
	}
	return rows, nil
}

func (r *Repository) Create(
	ctx context.Context,
	ticket *persistence.Ticket,
) (*persistence.Ticket, error) {

	q := query.Use(r.db)

	err := q.WithContext(ctx).Ticket.Create(ticket)
	if err != nil {
		return nil, fmt.Errorf("repo: create ticket %w", err)
	}

	return ticket, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*persistence.Ticket, error) {

	q := query.Use(r.db)
	ticket, err := q.WithContext(ctx).Ticket.Where(q.Ticket.ID.Eq(id)).First()
	if err != nil {
		return nil, fmt.Errorf("repo: get ticket %w", err)
	}

	return ticket, nil
}

func (r *Repository) Update(ctx context.Context, ticket *persistence.Ticket) error {
	q := query.Use(r.db)
	return q.WithContext(ctx).Ticket.Save(ticket)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	q := query.Use(r.db)
	_, err := q.WithContext(ctx).Ticket.Where(q.Ticket.ID.Eq(id)).Delete()
	return err
}
