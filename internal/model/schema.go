package model

import (
	"blessdarah/tuts/internal/lib"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// User model reflects the table in the database
type User struct {
	gorm.Model
	ID       *string `gorm:"primaryKey,column:id"`
	Name     string  `gorm:"column:name"`
	Email    string  `gorm:"column:email"`
	Password string  `gorm:"column:password"`
	Events   []Event `gorm:"foreignKey:UserID;references:ID"`
}

// Hooks/Interceptor to add user id durint create
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	hashedPassword, hashErr := lib.HashPassword(u.Password)
	if hashErr != nil {
		return fmt.Errorf("failed to hash password: %w", hashErr)
	}
	u.Password = hashedPassword
	u.ID = &id

	return
}

// EVent model reflects the table in the database
type Event struct {
	gorm.Model
	ID          string  `gorm:"primaryKey,column:id"`
	UserID      *string `gorm:"column:user_id;not null;index"`
	User        User    `gorm:"foreignKey:UserID;references:ID"`
	Name        string
	Description *string
	Venue       string
	StartDate   time.Time
	EndDate     time.Time
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New().String()
	return
}
