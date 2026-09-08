package ticket_test

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/event"
	"blessdarah/tuts/internal/ticket"
	"blessdarah/tuts/internal/user"
	"blessdarah/tuts/pkg"
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jose/go-jose/v4/testutils/assert"
	"gorm.io/gorm"
)

type ticketRepositoryTestCase struct {
	t         *testing.T
	db        *gorm.DB
	repo      *ticket.Repository
	eventRepo *event.Repository
	userRepo  *user.Repository
}

func NewTicketRepositoryTestCase(
	t *testing.T,
) *ticketRepositoryTestCase {
	t.Helper()
	db := pkg.NewPostgres(t)
	repo := ticket.NewRepository(db)
	eventRepo := event.NewRepository(db)
	userRepo := user.NewRepository(db)

	fakeuser := user.FakeUserPersistence()

	fakeEvent := event.FakeEventPersistence(gofakeit.UUID())
	fakeEvent.UserID = fakeuser.ID

	return &ticketRepositoryTestCase{
		t:         t,
		db:        db,
		repo:      repo,
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (r *ticketRepositoryTestCase) setup(u *persistence.User, e *persistence.Event) {

	// insert a fake user
	_, err := r.userRepo.Create(*u)
	assert.NoError(r.t, err, "user repo: create user")

	// insert a new event
	_, err = r.eventRepo.Create(context.TODO(), *e)
	assert.NoError(r.t, err, "event repo: create event")
}
func TestCreateTicket(t *testing.T) {
	deps := NewTicketRepositoryTestCase(t)
	fakeuser := user.FakeUserPersistence()

	fakeEvent := event.FakeEventPersistence(gofakeit.UUID())
	fakeEvent.UserID = fakeuser.ID

	fakeTicket := ticket.FakeTicketPersistence(gofakeit.UUID())
	fakeTicket.EventID = fakeEvent.ID

	t.Run("Should create a ticket", func(t *testing.T) {

		deps.setup(&fakeuser, fakeEvent)
		tickPers, err := deps.repo.Create(context.TODO(), fakeTicket)

		assert.NoError(t, err)

		assert.Equal(t, fakeTicket.ID, tickPers.ID)
		assert.Equal(t, fakeTicket.Type, tickPers.Type)
		assert.Equal(t, fakeTicket.Price, tickPers.Price)
		assert.Equal(t, fakeTicket.EventID, tickPers.EventID)
		assert.Equal(t, fakeTicket.Description, tickPers.Description)

		var count int
		stmt := deps.db.Raw("SELECT COUNT(*) FROM tickets")
		err = stmt.Scan(&count).Error

		assert.NoError(t, err, "repo: count tickets")
		assert.Equal(t, count, 1)
	})
}
