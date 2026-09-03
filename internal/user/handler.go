package user

import (
	"blessdarah/tuts/internal/config"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/model"
	"encoding/json"
	"errors"
	"log/slog"

	"net/http"
)

type userService interface {
	GetAll() []model.User
	AddUser(user model.User) (*string, error)
	GetByEmail(email string) (model.User, error)
}

type Handler struct {
	cfg    *config.AppEnv
	logger *slog.Logger
	app    userService
}

func NewHandler(cfg *config.AppEnv, app userService) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: config.NewLogger(cfg).With(slog.String("module", "user")),
		app:    app,
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("list users")
	users := h.app.GetAll()

	userList := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userList = append(userList, UserResponse{
			ID:    *user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	h.logger.Info("users fetched", "count", len(userList))
	lib.WriteJSON(w, r, http.StatusOK, userList)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var cr CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		h.logger.Error("decode create user", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "invalid JSON format",
		})
		return
	}

	err := cr.Validate()
	if len(err) > 0 {
		h.logger.Error("validate user", "error", err)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeValidationError,
			Title:  "Validation Failed",
			Status: http.StatusBadRequest,
			Detail: "one or more fields are invalid",
			Errors: err.Fields(),
		})
		return
	}

	u := cr.ToUser()
	_, getErr := h.app.GetByEmail(u.Email)
	if getErr == nil {
		h.logger.Error("duplicate user", "email", u.Email)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeDuplicateRes,
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: "user with this email already exists",
		})
		return
	}

	if !errors.Is(getErr, ErrUserNotFound) {
		h.logger.Error("get user by email", "error", getErr)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to verify existing user",
		})
		return
	}

	userID, addErr := h.app.AddUser(u)
	if addErr != nil {
		h.logger.Error("add user", "error", addErr)
		lib.WriteProblem(w, r, lib.ProblemDetails{
			Type:   lib.ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to create user",
		})
		return
	}

	ur := UserResponse{ID: *userID, Name: u.Name, Email: u.Email}
	lib.WriteJSON(w, r, http.StatusCreated, ur)
}
