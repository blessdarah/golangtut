package ticket

import (
	"blessdarah/tuts/internal/lib"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateRequest struct {
	Type        string  `json:"type" fake:"{word}"`
	Price       float64 `json:"price" fake:"{float64range:0.01,100.0}"`
	EventID     string  `json:"eventId" fake:"{uuid}"`
	Description *string `json:"description" fake:"{sentence}"`
}

func (r *CreateRequest) Validate() error {
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
