package auth

import (
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"blessdarah/tuts/internal/user"
	"errors"
	"fmt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDuplicateUser = errors.New("duplicate user")

type userRepository interface {
	FindByEmail(email string) (model.User, error)
	FindById(id string) (model.User, error)
	Create(user model.User) (*string, error)
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

func (s *Service) ValidateCredentials(email, password string) (string, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}

	if err := lib.CheckPassword(password, u.Password); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}

	return *u.ID, nil
}

func (s *Service) Signup(u model.User) (*string, error) {
	_, err := s.repo.FindByEmail(u.Email)
	if err == nil {
		return nil, ErrDuplicateUser
	}

	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	return s.repo.Create(u)
}

func (s *Service) GetByID(id string) (model.User, error) {
	return s.repo.FindById(id)
}
