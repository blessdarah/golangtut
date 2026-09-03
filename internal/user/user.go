package user

import (
	"blessdarah/tuts/internal/lib"
	"fmt"

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
