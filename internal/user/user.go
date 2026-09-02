package user

import (
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
	u.ID = &id
	return
}
