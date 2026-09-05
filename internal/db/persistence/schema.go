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
}

func (Event) TableName() string {
	return "events"
}
