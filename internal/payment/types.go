package payment

import (
	"blessdarah/tuts/internal/lib"
	"net/http"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateRequest struct {
	EventID  *string  `json:"eventId,omitempty"`
	TicketID *string  `json:"ticketId,omitempty"`
	Amount   *float64 `json:"-"`
	Quantity int      `json:"quantity"`
	Total    *float64 `json:"total"`
	Provider string   `json:"provider"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
}

func (res *CreateRequest) Validate() lib.HttpValidationError {
	errs := v.ValidateStruct(res,
		v.Field(&res.EventID, v.Required, v.Length(36, 0)),
		v.Field(&res.TicketID, v.Required, v.Length(36, 0)),
		v.Field(&res.Total, v.When(res.Total != nil, v.Min(0.0))),
		v.Field(&res.Provider, v.Required, v.Length(2, 0)),
		v.Field(&res.Name, v.Required, v.Length(2, 0)),
		v.Field(&res.Email, v.Required, is.Email),
	)
	if errs == nil {
		return nil
	}

	return lib.FormatError(errs.Error())
}

var AllowedFilters = map[string]string{
	"eventId":  "event_id",
	"ticketId": "ticket_id",
	"customer": "customer_name",
	"email":    "customer_email",
	"provider": "payment_provider",
}

func PaymentFilters(r *http.Request) map[string]string {
	queries := r.URL.Query()
	filters := make(map[string]string)
	for k, v := range AllowedFilters {
		if queries.Has(k) {
			filters[v] = queries.Get(k)
		}
	}
	return filters
}
