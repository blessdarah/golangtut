package user

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUserRequest) Validate() error {
	errs := v.ValidateStruct(u,
		v.Field(&u.Name, v.Required, v.Length(2, 0)),
		v.Field(&u.Email, v.Required, is.Email),
		v.Field(&u.Password, v.Required, v.Length(8, 0), is.Alphanumeric),
	)

	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

type UpdateUser struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (u *UpdateUser) Validate() error {
	errs := v.ValidateStruct(u,
		v.Field(&u.Name, v.When(u.Name != nil, v.Length(2, 0))),
		v.Field(&u.Email, v.When(u.Email != nil, is.Email)),
	)
	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *CreateUserRequest) ToUser() model.User {
	return model.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func FakeUserModel() model.User {
	return model.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 32),
	}
}

func FakeUserPersistence() persistence.User {
	return persistence.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 32),
	}
}
