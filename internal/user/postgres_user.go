package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInternal     = errors.New("internal error")
)

// Accept interfaces and return structs

type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
// It takes a gorm database as a parameter
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// list returns all users
func (r *Repository) List() []User {
	var users []User
	r.db.Find(&users)

	return users
}

// Create create a new user record
// @returns error if any
func (r *Repository) Create(user User) (*string, error) {
	return user.ID, r.db.Create(&user).Error
}

// FindByEmail finds a user by email
// @returns user if found, error otherwise
func (r *Repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

// FindById finds a user by id (uuid)
// @returns user if found, error otherwise
func (r *Repository) FindById(id string) (User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error

	// if user is not found

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

// Update updates a user record
// @returns error if any
func (r *Repository) Update(user User) error {
	return r.db.Save(&user).Error
}

// Delete deletes a user record
// @returns error if any
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&User{}, id).Error
}
