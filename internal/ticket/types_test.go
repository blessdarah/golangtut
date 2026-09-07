package ticket_test

import (
	"blessdarah/tuts/internal/ticket"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-jose/go-jose/v4/testutils/assert"
)

func TestCreateRequest(t *testing.T) {
	type TestCase struct {
		name     string
		input    ticket.CreateRequest
		expected string
	}

	testCases := []TestCase{
		{
			name: "throws no error for valid request",
			input: func() ticket.CreateRequest {
				var req ticket.CreateRequest
				gofakeit.Struct(&req)
				return req
			}(),
			expected: "",
		},
		{
			name: "throws an error for invalid event id",
			input: func() ticket.CreateRequest {
				var req ticket.CreateRequest
				gofakeit.Struct(&req)
				req.EventID = "invalid"
				return req
			}(),
			expected: `{"errors":{"eventId":"eventId the length must be no less than 36"}}`,
		},
		{
			name: "throw an error for price less that 0.0",
			input: func() ticket.CreateRequest {
				var req ticket.CreateRequest
				gofakeit.Struct(&req)
				req.Price = -0.01
				return req
			}(),
			expected: `{"errors":{"price":"price must be no less than 0"}}`,
		},
		{
			name: "throw and error for empty type",
			input: func() ticket.CreateRequest {
				var req ticket.CreateRequest
				gofakeit.Struct(&req)
				req.Type = ""
				return req
			}(),
			expected: `{"errors":{"type":"type cannot be blank"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.input.Validate()

			var got string
			if errs != nil {
				got = errs.Error()
			}

			assert.Equal(t, tc.expected, got)
		})
	}
}

func RunTypesTest(t *testing.T) {
	TestCreateRequest(t)
}

