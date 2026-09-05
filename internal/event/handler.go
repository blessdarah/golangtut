package event

import (
	"blessdarah/tuts/internal/auth"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"context"
	"encoding/json"
	"log/slog"

	"net/http"
)

type service interface {
	GetAll(ctx context.Context) ([]*model.Event, error)
	GetByUserID(ctx context.Context) ([]*model.Event, error)
	Create(ctx context.Context, event model.Event) (*model.Event, error)
	Get(ctx context.Context, id string) (*model.Event, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	service service
	logger  *slog.Logger
}

func NewHandler(service service, l *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  l,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("failed to get all events", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to get all events",
			Status:   http.StatusInternalServerError,
			Detail:   "internal server error",
			Type:     "https://httpstatuses.com/500",
			Instance: "/events",
		})
		return
	}

	res := make([]Response, len(events))
	for i, e := range events {
		res[i] = ToResponse(e)
	}

	lib.WriteJSON(w, r, http.StatusOK, res)
}

func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.logger.Error("missing auth context user")
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "unauthorized",
			Status:   http.StatusUnauthorized,
			Detail:   "invalid or missing access token",
			Type:     lib.ProblemTypeValidationError,
			Instance: r.URL.Path,
		})
		return
	}

	events, err := h.service.GetByUserID(ctx)
	if err != nil {
		h.logger.Error("failed to get all events", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to get all events",
			Status:   http.StatusInternalServerError,
			Detail:   "internal server error",
			Type:     "https://httpstatuses.com/500",
			Instance: "/events",
		})
		return
	}

	res := make([]Response, len(events))
	for i, e := range events {
		res[i] = ToResponse(e)
	}

	lib.WriteJSON(w, r, http.StatusOK, res)

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		h.logger.Error("missing auth context user")
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "unauthorized",
			Status:   http.StatusUnauthorized,
			Detail:   "invalid or missing access token",
			Type:     lib.ProblemTypeValidationError,
			Instance: r.URL.Path,
		})
		return
	}

	var req CreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to read event", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to read event",
			Status:   http.StatusBadRequest,
			Detail:   "invalid request body",
			Type:     lib.ProblemTypeValidationError,
			Instance: r.URL.Path,
		})
		return
	}

	vErrs := req.Validate(ctx)
	if vErrs != nil {
		h.logger.Error("failed to validate event", "error", vErrs)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to validate event",
			Status:   http.StatusBadRequest,
			Detail:   "invalid request body",
			Type:     lib.ProblemTypeValidationError,
			Instance: r.URL.Path,
			Errors:   vErrs.Fields(),
		})
		return
	}

	eventModel, err := req.ToEvent()
	if err != nil {
		h.logger.Error("failed to parse event dates", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to parse event dates",
			Status:   http.StatusBadRequest,
			Detail:   "startDate and endDate must use format YYYY-MM-DD",
			Type:     lib.ProblemTypeValidationError,
			Instance: r.URL.Path,
		})
		return
	}

	eventModel.UserID = &userID
	ev, err := h.service.Create(ctx, eventModel)
	if err != nil {
		h.logger.Error("failed to create event", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Title:    "failed to create event",
			Status:   http.StatusInternalServerError,
			Detail:   "internal server error",
			Type:     "https://httpstatuses.com/500",
			Instance: "/events",
		})
		return
	}

	lib.WriteJSON(w, r, http.StatusCreated, ToResponse(ev))
}

func ToResponse(e *model.Event) Response {
	return Response{
		ID:          *e.ID,
		UserID:      *e.UserID,
		Name:        e.Name,
		Venue:       e.Venue,
		StartDate:   *e.StartDate,
		EndDate:     *e.EndDate,
		Description: e.Description,
		CreatedAt:   *e.CreatedAt,
		UpdatedAt:   *e.UpdatedAt,
	}
}
