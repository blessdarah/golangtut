package auth

import (
	"blessdarah/tuts/internal/lib"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *LoginRequest) Validate() lib.HttpValidationError {
	errs := v.ValidateStruct(u,
		v.Field(&u.Email, v.Required, is.Email),
		v.Field(&u.Password, v.Required, v.Length(8, 0), is.Alphanumeric),
	)

	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

type LoginResponse struct {
	Token string `json:"token"`
}
