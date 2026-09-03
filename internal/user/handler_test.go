package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"blessdarah/tuts/internal/config"
	"blessdarah/tuts/internal/lib"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type mockUserService struct {
	getAllUsers []User
	getAllErr   error

	getByEmailUser User
	getByEmailErr  error

	addUserID  *string
	addUserErr error
}

func (m *mockUserService) GetAll() ([]User, error) {
	return m.getAllUsers, m.getAllErr
}

func (m *mockUserService) AddUser(user User) (*string, error) {
	return m.addUserID, m.addUserErr
}

func (m *mockUserService) GetByEmail(email string) (User, error) {
	return m.getByEmailUser, m.getByEmailErr
}

func newTestHandler(svc userService) *Handler {
	cfg := &config.AppEnv{Debug: true, Env: "test", AppHost: "localhost"}
	return NewHandler(cfg, svc)
}

func executeRequest(t *testing.T, handler http.HandlerFunc, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Method(method, path, handler)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res
}

func decodeProblem(t *testing.T, res *httptest.ResponseRecorder) lib.ProblemDetails {
	t.Helper()
	var p lib.ProblemDetails
	if err := json.Unmarshal(res.Body.Bytes(), &p); err != nil {
		t.Fatalf("failed to decode problem: %v", err)
	}
	return p
}

func TestGetUsersSuccessBareJSON(t *testing.T) {
	id := "u-1"
	h := newTestHandler(&mockUserService{getAllUsers: []User{{ID: &id, Name: "Ada", Email: "ada@example.com"}}})
	res := executeRequest(t, h.GetUsers, http.MethodGet, "/users", nil)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json, got %s", got)
	}

	var users []UserResponse
	if err := json.Unmarshal(res.Body.Bytes(), &users); err != nil {
		t.Fatalf("failed to decode users: %v", err)
	}
	if len(users) != 1 || users[0].ID != id {
		t.Fatalf("unexpected users payload: %+v", users)
	}
}

func TestGetUsersInternalErrorProblem(t *testing.T) {
	h := newTestHandler(&mockUserService{getAllErr: errors.New("db down")})
	res := executeRequest(t, h.GetUsers, http.MethodGet, "/users", nil)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", res.Code)
	}
	if got := res.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("expected application/problem+json, got %s", got)
	}

	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeInternalError || p.Instance != "/users" || p.RequestID == "" {
		t.Fatalf("unexpected problem details: %+v", p)
	}
}

func TestCreateUserInvalidJSONProblem(t *testing.T) {
	h := newTestHandler(&mockUserService{})
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", []byte("{"))

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeValidationError || p.Instance != "/users" || p.RequestID == "" {
		t.Fatalf("unexpected problem details: %+v", p)
	}
}

func TestCreateUserValidationProblemIncludesErrors(t *testing.T) {
	h := newTestHandler(&mockUserService{})
	body := []byte(`{"name":"a","email":"bad","password":"1"}`)
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", body)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeValidationError {
		t.Fatalf("expected validation type, got %s", p.Type)
	}
	if len(p.Errors) == 0 {
		t.Fatalf("expected validation errors extension, got %+v", p)
	}
}

func TestCreateUserDuplicateProblem(t *testing.T) {
	h := newTestHandler(&mockUserService{getByEmailUser: User{Email: "ada@example.com"}, getByEmailErr: nil})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", body)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", res.Code)
	}
	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeDuplicateRes {
		t.Fatalf("expected duplicate type, got %s", p.Type)
	}
}

func TestCreateUserInternalOnLookupError(t *testing.T) {
	h := newTestHandler(&mockUserService{getByEmailErr: errors.New("db unavailable")})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", body)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", res.Code)
	}
	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeInternalError {
		t.Fatalf("expected internal type, got %s", p.Type)
	}
}

func TestCreateUserInternalOnCreateError(t *testing.T) {
	h := newTestHandler(&mockUserService{getByEmailErr: ErrUserNotFound, addUserErr: errors.New("insert failed")})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", body)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", res.Code)
	}
	p := decodeProblem(t, res)
	if p.Type != lib.ProblemTypeInternalError {
		t.Fatalf("expected internal type, got %s", p.Type)
	}
}

func TestCreateUserSuccessBareJSON(t *testing.T) {
	id := "new-user"
	h := newTestHandler(&mockUserService{getByEmailErr: ErrUserNotFound, addUserID: &id})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := executeRequest(t, h.CreateUser, http.MethodPost, "/users", body)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.Code)
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json, got %s", got)
	}

	var payload map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload) != 3 {
		t.Fatalf("expected bare resource object, got %+v", payload)
	}
}
