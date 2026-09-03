package auth

import (
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"blessdarah/tuts/internal/user"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type authService interface {
	Signup(user model.User) (*string, error)
	GetByID(id string) (model.User, error)
}

type Handler struct {
	service authService
	oauth   OAuthServer
	logger  *slog.Logger
}

func NewAuthHandler(s authService, o OAuthServer, l *slog.Logger) *Handler {
	return &Handler{
		service: s,
		oauth:   o,
		logger:  l,
	}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("decode signup request", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid JSON format",
		})
		return
	}

	if err := req.Validate(); len(err) > 0 {
		h.logger.Error("validate signup request", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Validation Failed",
			Status: http.StatusBadRequest,
			Detail: "one or more fields are invalid",
			Errors: err.Fields(),
		})
		return
	}

	u := req.ToUser()
	userID, err := h.service.Signup(u)
	if err != nil {
		h.logger.Error("signup user", "error", err)
		if errors.Is(err, ErrDuplicateUser) {
			lib.WriteProblem(w, r, lib.ProblemDetails{
				Type:   lib.ProblemTypeDuplicateRes,
				Title:  "Conflict",
				Status: http.StatusConflict,
				Detail: "user with this email already exists",
			})
			return
		}

		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to create user",
		})
		return
	}

	resp := user.UserResponse{
		ID:    *userID,
		Name:  u.Name,
		Email: u.Email,
	}
	lib.WriteJSON(w, r, http.StatusCreated, resp)
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	if err := h.oauth.HandleTokenRequest(w, r); err != nil {
		h.logger.Error("issue oauth token", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Unauthorized",
			Status: http.StatusUnauthorized,
			Detail: "invalid token request",
		})
		return
	}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		h.logger.Error("missing auth context user")
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Unauthorized",
			Status: http.StatusUnauthorized,
			Detail: "invalid or missing access token",
		})
		return
	}

	u, err := h.service.GetByID(userID)
	if err != nil {
		h.logger.Error("lookup auth user", "error", err)
		if errors.Is(err, user.ErrUserNotFound) {
			lib.WriteProblem(w, r, lib.ProblemDetails{
				Type:   lib.ProblemTypeNotFound,
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "user not found",
			})
			return
		}

		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to fetch user profile",
		})
		return
	}

	resp := user.UserResponse{
		ID:    *u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
	lib.WriteJSON(w, r, http.StatusOK, resp)
}
