package model

import (
	"time"
)

type User struct {
	ID        *string    `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type Event struct {
	ID          *string    `json:"id"`
	UserID      *string    `json:"userId"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Venue       string     `json:"venue"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type Ticket struct {
	ID          *string    `json:"id"`
	Type        string     `json:"type"`
	Price       float64    `json:"price"`
	EventID     string     `json:"eventId"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type Payment struct {
	ID        *string    `json:"id"`
	EventID   string     `json:"eventId"`
	TicketID  string     `json:"ticketId"`
	Amount    float64    `json:"amount"`
	Quantity  int        `json:"quantity"`
	Total     float64    `json:"total"`
	Provider  string     `json:"provider"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
