package payment

import (
	"blessdarah/tuts/internal/db/persistence"
	"blessdarah/tuts/internal/model"
	"context"
	"log/slog"
	"uuid"
)

type repository interface {
	List(ctx context.Context, filters map[string]string) ([]*persistence.Payment, error)
	Create(ctx context.Context, payment *persistence.Payment) (*persistence.Payment, error)
}

type Service struct {
	repo   repository
	logger *slog.Logger
}

func NewService(repo repository, logger *slog.Logger) *Service {
	return &Service{
		repo,
		logger,
	}
}

// GetAll returns all payments
func (s *Service) GetAll(ctx context.Context, filters map[string]string) ([]model.Payment, error) {
	s.logger.Info("get payments", "filters", filters)
	rows, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	s.logger.Info("transforming payments...")
	payments := make([]model.Payment, len(rows))
	for i := range rows {
		payments[i] = *toDomainPayment(rows[i])
	}

	s.logger.Info("payments fetched", "count", len(payments))
	return payments, nil
}

func (s *Service) MakePayment(ctx context.Context, req CreateRequest) (*model.Payment, error) {

	res := toPersistencePayment(req)
	created, err := s.repo.Create(ctx, &res)
	if err != nil {
		return nil, err
	}

	domainModel := toDomainPayment(created)
	return domainModel, nil
}

func toPersistencePayment(req CreateRequest) persistence.Payment {
	return persistence.Payment{
		ID:            uuid.New().String(),
		EventID:       *req.EventID,
		TicketID:      *req.TicketID,
		Amount:        *req.Amount,
		Quantity:      req.Quantity,
		Total:         *req.Total,
		Provider:      req.Provider,
		CustomerName:  req.Name,
		CustomerEmail: req.Email,
	}
}

func toDomainPayment(p *persistence.Payment) *model.Payment {
	return &model.Payment{
		ID:            &p.ID,
		EventID:       p.EventID,
		TicketID:      p.TicketID,
		CustomerName:  p.CustomerName,
		CustomerEmail: p.CustomerEmail,
		Amount:        p.Amount,
		Quantity:      p.Quantity,
		Total:         p.Total,
		Provider:      p.Provider,
		CreatedAt:     &p.CreatedAt,
		UpdatedAt:     &p.UpdatedAt,
	}
}
