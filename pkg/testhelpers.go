package pkg

import (
	"blessdarah/tuts/internal/db/persistence"
	"fmt"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres creates a new postgres database for testing
func NewPostgres(t *testing.T) *gorm.DB {
	t.Helper()

	postgresConfig := embeddedpostgres.DefaultConfig().
		Username("postgres").
		Password("postgres").
		Database("testdb").
		Port(9876)

	pg := embeddedpostgres.NewDatabase(postgresConfig)

	if err := pg.Start(); err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}

	t.Cleanup(func() {
		pg.Stop()
	})

	dsn := fmt.Sprintf(
		"host=localhost port=9876 user=postgres password=postgres dbname=testdb sslmode=disable",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	db.AutoMigrate(
		&persistence.User{},
		&persistence.Event{},
		&persistence.Ticket{},
		&persistence.Payment{},
	)

	return db
}
