package ticket

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/model"
	"context"
	"time"

	"github.com/google/uuid"
)

type ticketRepository interface {
	List(ctx context.Context) ([]*persistence.Ticket, error)
	Create(ctx context.Context, ticket *persistence.Ticket) (*persistence.Ticket, error)
	GetByID(ctx context.Context, id string) (*persistence.Ticket, error)
	Update(ctx context.Context, ticket *persistence.Ticket) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo ticketRepository
}

func NewService(repo ticketRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]model.Ticket, error) {
	tickets, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// convert to domain model
	domainTickets := make([]model.Ticket, len(tickets))
	for i := range tickets {
		domainTickets[i] = toDomainTicket(*tickets[i])
	}

	return domainTickets, nil
}

func (s *Service) Create(ctx context.Context, ticket model.Ticket) (*model.Ticket, error) {
	t := toPersistenceTicket(ticket)
	created, err := s.repo.Create(ctx, &t)
	if err != nil {
		return nil, err
	}

	modelTicket := toDomainTicket(*created)
	return &modelTicket, nil
}

func toDomainTicket(t persistence.Ticket) model.Ticket {

	return model.Ticket{
		ID:          &t.ID,
		Type:        t.Type,
		Price:       t.Price,
		EventID:     t.EventID,
		Description: t.Description,
		CreatedAt:   &t.CreatedAt,
		UpdatedAt:   &t.UpdatedAt,
	}
}

func toPersistenceTicket(t model.Ticket) persistence.Ticket {
	if t.ID == nil {
		id := uuid.New().String()
		t.ID = &id
	}

	now := time.Now()

	if t.CreatedAt == nil {
		t.CreatedAt = &now
	}

	if t.UpdatedAt == nil {
		t.UpdatedAt = &now
	}

	return persistence.Ticket{
		ID:          *t.ID,
		Type:        t.Type,
		Price:       t.Price,
		EventID:     t.EventID,
		Description: t.Description,
		CreatedAt:   *t.CreatedAt,
		UpdatedAt:   *t.UpdatedAt,
	}
}
