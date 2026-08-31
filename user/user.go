package user

import "errors"

type Store interface {
	Add(user User) error
	FindByEmail(email string) (User, error)
	RemoveByEmail(email string) error
	GetAll() []User
}

type Address struct {
	Zip    string
	City   string
	Street string
}

type User struct {
	Address
	Name  string
	Age   int
	Email string
}

type userStore struct {
	users []User
}

func NewUserStore() userStore {
	return userStore{
		users: []User{},
	}
}

func (u *User) validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	if len(u.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if u.Age < 18 {
		return errors.New("age must be at least 18")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

func (u *userStore) Add(user User) error {
	err := user.validate()
	if err != nil {
		return err
	}

	u.users = append(u.users, user)
	return nil
}

func (u *userStore) FindByEmail(email string) (User, error) {
	for _, user := range u.users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (u *userStore) RemoveByEmail(email string) error {
	_, err := u.FindByEmail(email)

	return err
}

func (u *userStore) GetAll() []User {
	return u.users
}
