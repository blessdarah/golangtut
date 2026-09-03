package user

import "blessdarah/tuts/internal/model"

type userRepository interface {
	List() []model.User
	Create(user model.User) (*string, error)
	FindByEmail(email string) (model.User, error)
	FindById(id string) (model.User, error)
	Update(user model.User) error
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
	return a.repo.List()
}

// Create creates a new user
// by delegating to the repository
func (a *Service) AddUser(user model.User) (*string, error) {
	return a.repo.Create(user)
}

// GetByEmail finds a user by email
// by delegating to the repository
func (a *Service) GetByEmail(email string) (model.User, error) {
	return a.repo.FindByEmail(email)
}

// GetByID finds a user by id
// by delegating to the repository
func (a *Service) GetByID(id string) (model.User, error) {
	return a.repo.FindById(id)
}

// UpdateUser updates a user
// by delegating to the repository
func (a *Service) UpdateUser(user model.User) error {
	return a.repo.Update(user)
}

// DeleteUser deletes a user
// by delegating to the repository
func (a *Service) DeleteUser(id string) error {
	return a.repo.Delete(id)
}
