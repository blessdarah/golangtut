package ticket_test

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/ticket"
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jose/go-jose/v4/testutils/assert"
)

type mockTicketRepository struct {
	t    *testing.T
	data []*persistence.Ticket
}

func NewRepository(t *testing.T) *mockTicketRepository {
	return &mockTicketRepository{
		t:    t,
		data: nil,
	}
}

func (r *mockTicketRepository) List(ctx context.Context) ([]*persistence.Ticket, error) {
	return r.data, nil
}

func (r *mockTicketRepository) GetByID(
	ctx context.Context,
	id string,
) (*persistence.Ticket, error) {
	for _, t := range r.data {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, nil
}

func (r *mockTicketRepository) Update(
	ctx context.Context,
	ticket *persistence.Ticket,
) error {
	for i, t := range r.data {
		if t.ID == ticket.ID {
			r.data[i] = ticket
			return nil
		}
	}
	return nil
}

func (r *mockTicketRepository) Delete(
	ctx context.Context,
	id string,
) error {
	for i, t := range r.data {
		if t.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *mockTicketRepository) Create(
	ctx context.Context,
	ticket *persistence.Ticket,
) (*persistence.Ticket, error) {
	r.data = append(r.data, ticket)
	return ticket, nil
}

func TestGetAll(t *testing.T) {

	repo := NewRepository(t)
	service := ticket.NewService(repo)
	t.Run("Should return empty slice when no records", func(t *testing.T) {

		tickets, err := service.GetAll(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if len(tickets) != 0 {
			t.Errorf("expected no tickets, got %d", len(tickets))
		}

		assert.Len(t, repo.data, 0)
	})

	t.Run("Should return all records", func(t *testing.T) {
		desc := gofakeit.Sentence(1)
		repo.data = []*persistence.Ticket{
			{
				ID:          "1",
				Type:        "type",
				Price:       1.0,
				EventID:     "event-id",
				Description: &desc,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          "2",
				Type:        "type",
				Price:       1.0,
				EventID:     "event-id",
				Description: &desc,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		tickets, err := service.GetAll(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if len(tickets) != 2 {
			t.Errorf("expected 2 tickets, got %d", len(tickets))
		}

		assert.Len(t, repo.data, 2)
	})
}

func RunServiceTest(t *testing.T) {
	TestGetAll(t)
}
