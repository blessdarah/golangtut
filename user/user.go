package user

import (
	"errors"
)

type Store interface {
	Add(user User)
	FindByEmail(email string) (User, error)
	RemoveByEmail(email string) error
	GetAll() []User
}

type ValidationError map[string]string

type Address struct {
	Zip    *string `json:"zip"`
	City   string  `json:"city" validate:"required"`
	Street string  `json:"street" validate:"required"`
}

type User struct {
	Address `json:"address"`
	Name    string `json:"name" validate:"required"`
	Age     *int   `json:"age,omitempty"`
	Email   string `json:"email" validate:"required"`
}

type userStore struct {
	users []User
}

func NewUserStore() *userStore {
	return &userStore{
		users: []User{},
	}
}

func (u *User) Validate() *ValidationError {
	errs := ValidationError{}
	if u.Name == "" {
		errs["name"] = "name is required"
	}

	if len(u.Name) < 3 {
		errs["name"] = "name must be at least 3 characters"
	}

	if u.Age != nil && *u.Age < 18 {
		errs["age"] = "age must be at least 18"
	}

	if u.Email == "" {
		errs["email"] = "email is required"
	}

	if u.Zip != nil && len(*u.Zip) != 5 {
		errs["zip"] = "zip code must be 5 digits"
	}

	if u.Street == "" {
		errs["street"] = "street is required"
	}

	if u.City == "" {
		errs["city"] = "city is required"
	}

	if len(errs) > 0 {
		return &errs
	}

	return nil
}

func (u *userStore) Add(user User) {
	u.users = append(u.users, user)
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
