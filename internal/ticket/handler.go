package ticket

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/event"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type ticketService interface {
	GetAll(ctx context.Context) ([]model.Ticket, error)
	Create(ctx context.Context, ticket model.Ticket) (*model.Ticket, error)
	// GetByID(ctx context.Context, id string) (*model.Ticket, error)
	// Update(ctx context.Context, ticket model.Ticket) error
	// Delete(ctx context.Context, id string) error
}

type eventReader interface {
	Get(ctx context.Context, id string) (*model.Event, error)
}

type Handler struct {
	svc    ticketService
	er     eventReader
	logger *slog.Logger
}

func NewHandler(svc ticketService, er eventReader, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		er:     er,
		logger: logger,
	}
}

func (h *Handler) GetTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.svc.GetAll(r.Context())
	if err != nil {
		h.logger.Error("get tickets", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to get tickets",
		})
		return
	}

	h.logger.Info("tickets fetched", "count", len(tickets))
	lib.WriteJSON(w, r, http.StatusOK, tickets)
}

// Create creates a new ticket
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.UserIDFromContext(r.Context())

	if !ok {
		h.logger.Error("unathorizied", "error", "no user id in context")
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "unathorized",
		})
		return
	}

	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("decode create ticket", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid JSON format",
		})
		return
	}

	valErrs := req.Validate()
	if valErrs != nil {
		h.logger.Error("validate ticket", "error", valErrs)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Validation Failed",
			Status: http.StatusBadRequest,
			Detail: "one or more fields are invalid",
			Errors: valErrs.Fields(),
		})
		return
	}

	_, findErr := h.er.Get(r.Context(), req.EventID)
	if findErr != nil {
		h.logger.Error("get event", "error", findErr)

		if errors.Is(findErr, event.ErrEventNotFound) {
			lib.WriteProblem(w, r, lib.ProblemDetails{
				Type:   lib.ProblemTypeNotFound,
				Title:  "Event Not Found",
				Status: http.StatusNotFound,
				Detail: fmt.Sprintf("event with id %s not found", req.EventID),
			})
			return
		}

		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to get event",
		})
		return
	}

	modelTicket := model.Ticket{
		Type:        req.Type,
		Price:       req.Price,
		EventID:     req.EventID,
		Description: req.Description,
	}

	ticket, err := h.svc.Create(r.Context(), modelTicket)
	if err != nil {
		h.logger.Error("create ticket", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to create ticket",
		})
		return
	}

	lib.WriteJSON(w, r, http.StatusCreated, ticket)
}
