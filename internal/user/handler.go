package user

import (
	"blessdarah/tuts/internal/config"
	"encoding/json"
	"log/slog"

	"net/http"
)

type userService interface {
	GetAll() []User
	AddUser(user User) (*string, error)
	GetByEmail(email string) (User, error)
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
	h.logger.Info("create user entry")
	users := h.app.GetAll()
	userList := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userList = append(userList, UserResponse{
			ID:    *user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	h.logger.Info("user created")
	json.NewEncoder(w).Encode(userList)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var cr CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		h.logger.Error("decode create user", "error", err)

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid json format",
		})
		return
	}

	err := cr.Validate()

	if len(err) > 0 {
		h.logger.Error("validate user", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(err)
		return
	}

	user := cr.ToUser()

	// get user by email
	_, getErr := h.app.GetByEmail(user.Email)
	if getErr == nil {
		h.logger.Error("duplicate user", "error", getErr)
		w.WriteHeader(http.StatusConflict)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "user with this email already exists",
		})
		return
	}

	userID, addErr := h.app.AddUser(user)
	if addErr != nil {
		h.logger.Error("add user", "error", addErr)
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}

	w.WriteHeader(201)
	ur := UserResponse{
		ID:    *userID,
		Name:  user.Name,
		Email: user.Email,
	}

	json.NewEncoder(w).Encode(ur)

}
