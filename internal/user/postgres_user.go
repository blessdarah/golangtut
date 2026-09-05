package user

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/db/query"
	"context"
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
func (r *Repository) List() []persistence.User {
	q := query.Use(r.db)
	rows, err := q.WithContext(context.Background()).User.Find()
	if err != nil {
		return []persistence.User{}
	}

	users := make([]persistence.User, len(rows))
	for i, row := range rows {
		users[i] = *row
	}

	return users
}

// Create create a new user record
// @returns error if any
func (r *Repository) Create(user persistence.User) (*persistence.User, error) {
	q := query.Use(r.db)
	if err := q.WithContext(context.Background()).User.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
// @returns user if found, error otherwise
func (r *Repository) FindByEmail(email string) (persistence.User, error) {
	q := query.Use(r.db)
	user, err := q.WithContext(context.Background()).User.Where(q.User.Email.Eq(email)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return persistence.User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if err != nil {
		return persistence.User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return *user, nil
}

// FindById finds a user by id (uuid)
// @returns user if found, error otherwise
func (r *Repository) FindById(id string) (persistence.User, error) {
	q := query.Use(r.db)
	user, err := q.WithContext(context.Background()).User.Where(q.User.ID.Eq(id)).First()

	// if user is not found

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return persistence.User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if err != nil {
		return persistence.User{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return *user, nil
}

// Update updates a user record
// @returns error if any
func (r *Repository) Update(user persistence.User) error {
	q := query.Use(r.db)
	return q.WithContext(context.Background()).User.Save(&user)
}

// Delete deletes a user record
// @returns error if any
func (r *Repository) Delete(id string) error {
	q := query.Use(r.db)
	_, err := q.WithContext(context.Background()).User.Where(q.User.ID.Eq(id)).Delete()
	return err
}
