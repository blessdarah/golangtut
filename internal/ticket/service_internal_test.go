package ticket

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/model"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jose/go-jose/v4/testutils/assert"
)

func TestToDomainTicket(t *testing.T) {
	t.Run("Should convert persistence to domain ticket", func(t *testing.T) {
		id := gofakeit.UUID()
		type_ := gofakeit.Word()
		price := gofakeit.Float64Range(0.0, 100.0)
		eventID := gofakeit.UUID()
		desc := gofakeit.Sentence(1)
		now := gofakeit.DateRange(time.Now(), time.Now().Add(time.Hour*24))
		createdAt := now
		updatedAt := now

		ticket := persistence.Ticket{
			ID:          id,
			Type:        type_,
			Price:       price,
			EventID:     eventID,
			Description: &desc,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		got := toDomainTicket(ticket)

		assert.Equal(t, ticket.ID, *got.ID)
		assert.Equal(t, ticket.Type, got.Type)
		assert.Equal(t, ticket.Price, got.Price)
		assert.Equal(t, ticket.EventID, got.EventID)
		assert.Equal(t, ticket.Description, got.Description)
	})
}

func TestToPersistenceTicket(t *testing.T) {
	id := gofakeit.UUID()
	des := gofakeit.Sentence(1)
	domainModel := model.Ticket{
		ID:          &id,
		Type:        gofakeit.Word(),
		Price:       gofakeit.Float64Range(0.0, 100.0),
		EventID:     gofakeit.UUID(),
		Description: &des,
	}

	got := toPersistenceTicket(domainModel)

	assert.Equal(t, *domainModel.ID, got.ID)
	assert.Equal(t, domainModel.Type, got.Type)
	assert.Equal(t, domainModel.Price, got.Price)
	assert.Equal(t, domainModel.EventID, got.EventID)
	assert.Equal(t, domainModel.Description, got.Description)
}
