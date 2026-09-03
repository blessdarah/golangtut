package event

import (
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateRequest struct {
	Name        string  `json:"name"`
	Venue       string  `json:"venue"`
	StartDate   string  `json:"startDate"`
	EndDate     string  `json:"endDate"`
	Description *string `json:"description"`
}

func (res *CreateRequest) Validate(ctx context.Context) lib.HttpValidationError {
	errs := validation.ValidateStructWithContext(ctx, res,
		validation.Field(&res.Name, validation.Required, validation.Length(2, 0)),
		validation.Field(&res.Venue, validation.Required, validation.Length(2, 0)),
		validation.Field(&res.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&res.EndDate,
			validation.Required,
			validation.Date("2006-01-02"),
			validation.By(func(value any) error {
				startDate, startErr := time.Parse("2006-01-02", res.StartDate)
				endDate, endErr := time.Parse("2006-01-02", res.EndDate)
				if startErr != nil || endErr != nil {
					return nil
				}

				if !startDate.Before(endDate) {
					return errors.New("must be after startDate")
				}

				return nil
			}),
		),
		validation.Field(&res.Description, validation.When(res.Description != nil, validation.Length(2, 0))),
	)

	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

func (res *CreateRequest) ToEvent() (model.Event, error) {
	startDate, err := time.Parse("2006-01-02", res.StartDate)
	if err != nil {
		return model.Event{}, err
	}

	endDate, err := time.Parse("2006-01-02", res.EndDate)
	if err != nil {
		return model.Event{}, err
	}

	return model.Event{
		Name:        res.Name,
		Venue:       res.Venue,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: res.Description,
	}, nil
}

type UpdateResponse struct {
	Name        *string `json:"name"`
	Venue       *string `json:"venue"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Description *string `json:"description"`
}

func (res *UpdateResponse) Validate(ctx context.Context) lib.HttpValidationError {
	errs := validation.ValidateStructWithContext(ctx, res,
		validation.Field(&res.Name, validation.When(res.Name != nil, validation.Length(2, 0))),
		validation.Field(&res.Venue, validation.When(res.Venue != nil, validation.Length(2, 0))),
		validation.Field(&res.StartDate, validation.When(res.StartDate != nil, validation.Date("2006-01-02"))),
		validation.Field(&res.EndDate,
			validation.When(res.EndDate != nil, validation.Date("2006-01-02")),
			validation.By(func(value interface{}) error {
				if res.StartDate == nil || res.EndDate == nil {
					return nil
				}

				startDate, startErr := time.Parse("2006-01-02", *res.StartDate)
				endDate, endErr := time.Parse("2006-01-02", *res.EndDate)
				if startErr != nil || endErr != nil {
					return nil
				}

				if !startDate.Before(endDate) {
					return errors.New("must be after startDate")
				}

				return nil
			}),
		),
		validation.Field(&res.Description, validation.When(res.Description != nil, validation.Length(2, 0))),
	)

	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

type Response struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Venue       string    `json:"venue"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
