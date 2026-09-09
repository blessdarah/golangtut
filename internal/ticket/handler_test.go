package ticket_test

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/event"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"blessdarah/tuts/internal/ticket"
	"blessdarah/tuts/internal/user"
	"blessdarah/tuts/pkg"

	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

const (
	oauthClientID     = "e2e-client-id"
	oauthClientSecret = "e2e-client-secret"
)

type e2eTicket struct {
	router  http.Handler
	db      *gorm.DB
	token   string
	eventID string
}

func newE2ETicket(t *testing.T) *e2eTicket {
	t.Helper()

	db := pkg.NewPostgres(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	userSvc := user.NewService(user.NewRepository(db))
	authSvc := auth.NewService(userSvc)
	oauthServer, err := auth.NewOAuthServer(oauthClientID, oauthClientSecret, 60, 24, authSvc)
	if err != nil {
		t.Fatalf("new oauth server: %v", err)
	}
	authMw := auth.RequireBearer(oauthServer, logger)

	eventSvc := event.NewService(event.NewRepository(db))
	ticketHandler := ticket.NewHandler(ticket.NewService(ticket.NewRepository(db)), eventSvc, logger)

	userID, err := userSvc.AddUser(model.User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, false, false, 12),
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	token, err := oauthServer.IssueTestToken(*userID)
	if err != nil {
		t.Fatalf("issue test token: %v", err)
	}

	fakeEventModel := event.FakeEventModel(*userID)
	ev, err := eventSvc.Create(context.TODO(), *fakeEventModel)
	if err != nil {
		t.Fatalf("create event: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Get("/tickets", ticketHandler.GetTickets)
	r.Group(func(gr chi.Router) {
		gr.Use(authMw)
		gr.Post("/tickets", ticketHandler.Create)
		gr.Get("/events/{id}/{ticket_id}/pay", ticketHandler.GetTicket)
	})

	return &e2eTicket{
		router:  r,
		db:      db,
		token:   token,
		eventID: *ev.ID,
	}
}

func (e *e2eTicket) do(t *testing.T, method, path, token, contentType, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)
	return rec
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("decode response body %q: %v", rec.Body.String(), err)
	}
}

func decodeProblem(t *testing.T, rec *httptest.ResponseRecorder) lib.ProblemDetails {
	t.Helper()
	var problem lib.ProblemDetails
	decodeJSON(t, rec, &problem)
	return problem
}

func assertProblem(t *testing.T, rec *httptest.ResponseRecorder, status int, title, detail string) {
	t.Helper()
	if rec.Code != status {
		t.Errorf("status code = %d, want %d, body = %s", rec.Code, status, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("content type = %q, want %q", got, "application/problem+json")
	}
	problem := decodeProblem(t, rec)
	if problem.Title != title {
		t.Errorf("title = %q, want %q", problem.Title, title)
	}
	if problem.Detail != detail {
		t.Errorf("detail = %q, want %q", problem.Detail, detail)
	}
}

func TestTicketAPIEndToEnd(t *testing.T) {
	env := newE2ETicket(t)

	var ticketID string

	t.Run("create ticket returns the created ticket", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"vip","price":49.99,"eventId":%q,"description":"vip pass"}`, env.eventID)
		rec := env.do(t, http.MethodPost, "/tickets", env.token, "application/json", body)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status code = %d, body = %s", rec.Code, rec.Body.String())
		}

		var resp model.Ticket
		decodeJSON(t, rec, &resp)
		if resp.ID == nil || *resp.ID == "" {
			t.Fatal("ticket id is empty")
		}
		if resp.Type != "vip" {
			t.Errorf("type = %q, want %q", resp.Type, "vip")
		}
		if resp.Price != 49.99 {
			t.Errorf("price = %v, want %v", resp.Price, 49.99)
		}
		if resp.EventID != env.eventID {
			t.Errorf("eventId = %q, want %q", resp.EventID, env.eventID)
		}
		if resp.Description == nil || *resp.Description != "vip pass" {
			t.Errorf("description = %v, want %q", resp.Description, "vip pass")
		}
		ticketID = *resp.ID
	})

	t.Run("lists all tickets in db", func(t *testing.T) {
		rec := env.do(t, http.MethodGet, "/tickets", "", "", "")

		if rec.Code != http.StatusOK {
			t.Fatalf("status code = %d, body = %s", rec.Code, rec.Body.String())
		}

		var tickets []model.Ticket
		decodeJSON(t, rec, &tickets)

		if len(tickets) == 0 {
			t.Fatal("no tickets found")
		}
	})

	t.Run("create ticket without a token returns unauthorized", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"regular","price":19.99,"eventId":%q}`, env.eventID)
		rec := env.do(t, http.MethodPost, "/tickets", "", "application/json", body)

		assertProblem(t, rec, http.StatusUnauthorized, "Unauthorized", "invalid or missing access token")
	})

	t.Run("create ticket with an invalid token returns unauthorized", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"regular","price":19.99,"eventId":%q}`, env.eventID)
		rec := env.do(t, http.MethodPost, "/tickets", "not-a-real-token", "application/json", body)

		assertProblem(t, rec, http.StatusUnauthorized, "Unauthorized", "invalid or missing access token")
	})

	t.Run("create ticket with malformed json returns bad request", func(t *testing.T) {
		rec := env.do(t, http.MethodPost, "/tickets", env.token, "application/json", `{"type":`)

		assertProblem(t, rec, http.StatusBadRequest, "Bad Request", "invalid JSON format")
	})

	t.Run("create ticket with an invalid price returns validation errors", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"regular","price":-0.01,"eventId":%q}`, env.eventID)
		rec := env.do(t, http.MethodPost, "/tickets", env.token, "application/json", body)

		assertProblem(t, rec, http.StatusBadRequest, "Validation Failed", "one or more fields are invalid")
		problem := decodeProblem(t, rec)
		wantErrs := map[string]string{"price": "price must be no less than 0"}
		if !reflect.DeepEqual(problem.Errors, wantErrs) {
			t.Errorf("errors = %v, want %v", problem.Errors, wantErrs)
		}
	})

	t.Run("create ticket for a missing event returns not found", func(t *testing.T) {
		missingEventID := gofakeit.UUID()
		body := fmt.Sprintf(`{"type":"regular","price":19.99,"eventId":%q}`, missingEventID)
		rec := env.do(t, http.MethodPost, "/tickets", env.token, "application/json", body)

		assertProblem(t, rec, http.StatusNotFound, "Event Not Found", fmt.Sprintf("event with id %s not found", missingEventID))
	})

	t.Run("list tickets includes the created ticket", func(t *testing.T) {
		rec := env.do(t, http.MethodGet, "/tickets", "", "", "")

		if rec.Code != http.StatusOK {
			t.Fatalf("status code = %d, body = %s", rec.Code, rec.Body.String())
		}

		var tickets []model.Ticket
		decodeJSON(t, rec, &tickets)

		var found *model.Ticket
		for i := range tickets {
			if tickets[i].ID != nil && *tickets[i].ID == ticketID {
				found = &tickets[i]
				break
			}
		}
		if found == nil {
			t.Fatalf("ticket %s not found in list, body = %s", ticketID, rec.Body.String())
		}
		if found.Type != "vip" || found.Price != 49.99 || found.EventID != env.eventID {
			t.Errorf("listed ticket = %+v, want type vip, price 49.99, eventId %s", *found, env.eventID)
		}
	})

	t.Run("ticket endpoints return internal server errors when the tickets table is unavailable", func(t *testing.T) {
		if err := env.db.Exec("DROP TABLE tickets CASCADE").Error; err != nil {
			t.Fatalf("drop tickets table: %v", err)
		}

		rec := env.do(t, http.MethodGet, "/tickets", "", "", "")
		assertProblem(t, rec, http.StatusInternalServerError, "Internal Server Error", "failed to get tickets")

		body := fmt.Sprintf(`{"type":"regular","price":19.99,"eventId":%q}`, env.eventID)
		rec = env.do(t, http.MethodPost, "/tickets", env.token, "application/json", body)
		assertProblem(t, rec, http.StatusInternalServerError, "Internal Server Error", "failed to create ticket")
	})
}
