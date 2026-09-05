package ticket

import (
	"blessdarah/tuts/internal/lib"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateRequest struct {
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	EventID     string  `json:"eventId"`
	Description *string `json:"description"`
}

func (r *CreateRequest) Validate() lib.HttpValidationError {
	errs := validation.ValidateStruct(r,
		validation.Field(&r.Type, validation.Required, validation.Length(2, 0)),
		validation.Field(&r.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&r.EventID, validation.Required, validation.Length(36, 0)),
		validation.Field(&r.Description, validation.When(r.Description != nil, validation.Length(0, 1000))),
	)

	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}
