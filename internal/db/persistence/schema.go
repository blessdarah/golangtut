package persistence

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"column:id;type:text;primaryKey"`
	Name      string         `gorm:"column:name;type:text;not null"`
	Email     string         `gorm:"column:email;type:text;not null;unique"`
	Password  string         `gorm:"column:password;type:text;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

type Event struct {
	ID          string         `gorm:"column:id;type:text;primaryKey"`
	UserID      string         `gorm:"column:user_id;type:text;not null;index:events_user_id_idx"`
	User        User           `gorm:"foreignKey:UserID;references:ID"`
	Name        string         `gorm:"column:name;type:text;not null"`
	Description *string        `gorm:"column:description;type:text"`
	Venue       string         `gorm:"column:venue;type:text;not null"`
	StartDate   time.Time      `gorm:"column:start_date;not null;index:events_start_date_idx"`
	EndDate     time.Time      `gorm:"column:end_date;not null;index:events_end_date_idx"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
	Payments    []Payment      `gorm:"foreignKey:EventID;references:ID"`
}

func (Event) TableName() string {
	return "events"
}

type Ticket struct {
	ID          string         `gorm:"column:ticket_id;type:text;primaryKey"`
	Type        string         `gorm:"column:type;type:varchar(30);not null"`
	Price       float64        `gorm:"column:price;type:float;not null"`
	EventID     string         `gorm:"column:event_id;type:text;not null;index:tickets_event_id_idx"`
	Event       Event          `gorm:"foreignKey:EventID;references:ID"`
	Description *string        `gorm:"column:description;type:text"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Payment struct {
	ID        string         `gorm:"column:id;type:text;primaryKey"`
	EventID   string         `gorm:"column:event_id;type:text;not null;index:payments_event_id_idx"`
	Event     Event          `gorm:"foreignKey:EventID;references:ID"`
	TicketID  string         `gorm:"column:ticket_id;type:text;not null;index:payments_ticket_id_idx"`
	Ticket    Ticket         `gorm:"foreignKey:TicketID;references:ID"`
	Amount    float64        `gorm:"column:amount;type:float;not null"`
	Quantity  int            `gorm:"column:quantity;type:integer;not null"`
	Total     float64        `gorm:"column:total;type:float;not null"`
	Provider  string         `gorm:"column:payment_provider;type:varchar(30);not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
