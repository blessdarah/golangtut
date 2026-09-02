package user

type userRepository interface {
	List() []User
	Create(user User) (*string, error)
	FindByEmail(email string) (User, error)
	FindById(id string) (User, error)
	Update(user User) error
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
func (a *Service) GetAll() []User {
	return a.repo.List()
}

// Create creates a new user
// by delegating to the repository
func (a *Service) AddUser(user User) (*string, error) {
	return a.repo.Create(user)
}

// GetByEmail finds a user by email
// by delegating to the repository
func (a *Service) GetByEmail(email string) (User, error) {
	return a.repo.FindByEmail(email)
}

// GetByID finds a user by id
// by delegating to the repository
func (a *Service) GetByID(id string) (User, error) {
	return a.repo.FindById(id)
}

// UpdateUser updates a user
// by delegating to the repository
func (a *Service) UpdateUser(user User) error {
	return a.repo.Update(user)
}

// DeleteUser deletes a user
// by delegating to the repository
func (a *Service) DeleteUser(id string) error {
	return a.repo.Delete(id)
}
