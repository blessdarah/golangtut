package user

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"fmt"

	"github.com/google/uuid"
)

type userRepository interface {
	List() []persistence.User
	Create(user persistence.User) (*persistence.User, error)
	FindByEmail(email string) (persistence.User, error)
	FindById(id string) (persistence.User, error)
	Update(user persistence.User) error
	Delete(id string) error
}

type Service struct {
	repo userRepository
}

// NewService creates a new app for user
func NewService(repo userRepository) *Service {
	return &Service{
		repo,
	}
}

// GetAll returns all users
// by delegating to the repository
func (a *Service) GetAll() []model.User {
	rows := a.repo.List()
	users := make([]model.User, len(rows))
	for i := range rows {
		users[i] = toDomainUser(rows[i])
	}

	return users
}

// Create creates a new user
// by delegating to the repository
func (a *Service) AddUser(user model.User) (*string, error) {
	id := uuid.NewString()
	hashedPassword, hashErr := lib.HashPassword(user.Password)
	if hashErr != nil {
		return nil, fmt.Errorf("%w: failed to hash password", ErrInternal)
	}

	user.ID = &id
	user.Password = hashedPassword
	created, err := a.repo.Create(toPersistenceUser(user))
	if err != nil {
		return nil, err
	}

	return &created.ID, nil
}

// GetByEmail finds a user by email
// by delegating to the repository
func (a *Service) GetByEmail(email string) (model.User, error) {
	u, err := a.repo.FindByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u), nil
}

// GetByID finds a user by id
// by delegating to the repository
func (a *Service) GetByID(id string) (model.User, error) {
	u, err := a.repo.FindById(id)
	if err != nil {
		return model.User{}, err
	}

	return toDomainUser(u), nil
}

func (a *Service) FindByEmail(email string) (model.User, error) {
	return a.GetByEmail(email)
}

func (a *Service) FindById(id string) (model.User, error) {
	return a.GetByID(id)
}

func (a *Service) Create(user model.User) (*string, error) {
	return a.AddUser(user)
}

// UpdateUser updates a user
// by delegating to the repository
func (a *Service) UpdateUser(user model.User) error {
	return a.repo.Update(toPersistenceUser(user))
}

// DeleteUser deletes a user
// by delegating to the repository
func (a *Service) DeleteUser(id string) error {
	return a.repo.Delete(id)
}

func toDomainUser(u persistence.User) model.User {
	id := u.ID
	return model.User{
		ID:        &id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toPersistenceUser(u model.User) persistence.User {
	id := ""
	if u.ID != nil {
		id = *u.ID
	}

	return persistence.User{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
