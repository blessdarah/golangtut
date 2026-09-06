package payment

import (
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"blessdarah/tuts/pkg"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type paymentService interface {
	MakePayment(ctx context.Context, req CreateRequest) (*model.Payment, error)
	GetAll(ctx context.Context, filters map[string]string) ([]model.Payment, error)
}

type eventReader interface {
	Get(ctx context.Context, id string) (*model.Event, error)
}

type ticketReader interface {
	GetByID(ctx context.Context, id string) (*model.Ticket, error)
}

type Handler struct {
	svc    paymentService
	logger *slog.Logger
	exr    pkg.ExchangeRateService
	er     eventReader
	tr     ticketReader
}

func NewHandler(
	svc paymentService,
	logger *slog.Logger,
	exr pkg.ExchangeRateService,
	er eventReader,
	tr ticketReader,
) *Handler {
	return &Handler{
		svc,
		logger,
		exr,
		er,
		tr,
	}
}

func (h *Handler) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.svc.GetAll(r.Context(), PaymentFilters(r))
	if err != nil {
		h.logger.Error("get payments", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to get payments",
		})
		return
	}

	h.logger.Info("payments fetched", "count", len(payments))
	lib.WriteJSON(w, r, http.StatusOK, payments)
}

// events/{id}/{ticket_id}/pay
func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {

	eventID := chi.URLParam(r, "id")
	ticketID := chi.URLParam(r, "ticket_id")
	if _, err := uuid.Parse(eventID); err != nil {
		h.logger.Error("parse event id", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid event id, exepected a uuid",
		})
		return
	}

	if _, err := uuid.Parse(ticketID); err != nil {
		h.logger.Error("parse ticket id", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid ticket id, exepected a uuid",
		})
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("decode create payment", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid JSON format",
		})
		return
	}
	req.EventID = &eventID
	req.TicketID = &ticketID

	vErrs := req.Validate()
	if vErrs != nil {
		h.logger.Error("validate payment", "error", vErrs)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Validation Failed",
			Status: http.StatusBadRequest,
			Detail: "one or more fields are invalid",
			Errors: vErrs.Fields(),
		})
		return
	}

	// check if event exists with id
	_, err := h.er.Get(r.Context(), *req.EventID)
	if err != nil {
		h.logger.Error("get event", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeNotFound,
			Title:  "Event Not Found",
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("event with id: %s not found", *req.EventID),
		})
		return
	}

	// check if ticket exists with id
	t, err := h.tr.GetByID(r.Context(), *req.TicketID)
	if err != nil {
		h.logger.Error("get ticket", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeNotFound,
			Title:  "Ticket Not Found",
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("ticket with id: %s not found", *req.TicketID),
		})
		return
	}

	exrRes, err := h.exr.GetExchangeRateByPair("USD", "XAF")
	if err != nil {
		h.logger.Error("get exchange rate", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to get exchange rate",
		})
		return
	}

	amt := t.Price * exrRes.ConversionRate
	req.Amount = &amt
	total := amt * float64(req.Quantity)
	req.Total = &total

	payment, err := h.svc.MakePayment(r.Context(), req)
	if err != nil {
		h.logger.Error("create payment", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to create payment",
		})
		return
	}

	lib.WriteJSON(w, r, http.StatusCreated, payment)
}
