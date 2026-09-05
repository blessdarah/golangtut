package model

import (
	"time"
)

type User struct {
	ID        *string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Event struct {
	ID          string
	UserID      *string
	Name        string
	Description *string
	Venue       string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
