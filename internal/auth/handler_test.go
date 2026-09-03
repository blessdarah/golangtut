package auth

import (
	"blessdarah/tuts/internal/config"
	"blessdarah/tuts/internal/lib"
	"blessdarah/tuts/internal/user"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type mockAuthService struct {
	signupID  *string
	signupErr error

	getByIDUser user.User
	getByIDErr  error
}

func (m *mockAuthService) Signup(u user.User) (*string, error) {
	return m.signupID, m.signupErr
}

func (m *mockAuthService) GetByID(id string) (user.User, error) {
	return m.getByIDUser, m.getByIDErr
}

type mockOAuthServer struct {
	handleErr error
	userID    string
	validErr  error
}

func (m *mockOAuthServer) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	if m.handleErr != nil {
		return m.handleErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"access_token":"abc","token_type":"Bearer","expires_in":3600}`))
	return nil
}

func (m *mockOAuthServer) ValidateBearerToken(r *http.Request) (string, error) {
	if m.validErr != nil {
		return "", m.validErr
	}
	return m.userID, nil
}

func newTestHandler(s authService, o OAuthServer) *Handler {
	cfg := &config.AppEnv{Debug: true, Env: "test", AppHost: "localhost"}
	logger := config.NewLogger(cfg)
	return NewAuthHandler(s, o, logger)
}

func doRequest(t *testing.T, method, path string, body []byte, h http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Method(method, path, h)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	return res
}

func TestSignupSuccess(t *testing.T) {
	id := "user-1"
	h := newTestHandler(&mockAuthService{signupID: &id}, &mockOAuthServer{})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := doRequest(t, http.MethodPost, "/auth/signup", body, h.Signup)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	if _, ok := got["password"]; ok {
		t.Fatalf("password must not be present in signup response")
	}
}

func TestSignupDuplicate(t *testing.T) {
	h := newTestHandler(&mockAuthService{signupErr: ErrDuplicateUser}, &mockOAuthServer{})
	body := []byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)
	res := doRequest(t, http.MethodPost, "/auth/signup", body, h.Signup)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", res.Code)
	}

	var problem lib.ProblemDetails
	if err := json.Unmarshal(res.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to parse problem: %v", err)
	}
	if problem.Type != lib.ProblemTypeDuplicateRes {
		t.Fatalf("expected duplicate resource type, got %s", problem.Type)
	}
}

func TestTokenSuccess(t *testing.T) {
	h := newTestHandler(&mockAuthService{}, &mockOAuthServer{})
	res := doRequest(t, http.MethodPost, "/oauth/token", nil, h.Token)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
}

func TestTokenFailure(t *testing.T) {
	h := newTestHandler(&mockAuthService{}, &mockOAuthServer{handleErr: errors.New("bad request")})
	res := doRequest(t, http.MethodPost, "/oauth/token", nil, h.Token)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}

func TestMeUnauthorized(t *testing.T) {
	h := newTestHandler(&mockAuthService{}, &mockOAuthServer{})
	res := doRequest(t, http.MethodGet, "/auth/me", nil, h.Me)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}

func TestMeSuccess(t *testing.T) {
	id := "user-1"
	h := newTestHandler(&mockAuthService{getByIDUser: user.User{ID: &id, Name: "Ada", Email: "ada@example.com"}}, &mockOAuthServer{userID: id})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(RequireBearer(&mockOAuthServer{userID: id}, nil))
	r.Get("/auth/me", h.Me)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if _, ok := got["password"]; ok {
		t.Fatalf("password must not be present")
	}
}

func TestRequireBearerBlocksInvalidToken(t *testing.T) {
	mw := RequireBearer(&mockOAuthServer{validErr: errors.New("bad token")}, nil)

	called := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if called {
		t.Fatalf("handler should not be called for invalid token")
	}
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}

func TestRequireBearerPassesUserContext(t *testing.T) {
	const id = "user-123"
	mw := RequireBearer(&mockOAuthServer{userID: id}, nil)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID, ok := userIDFromContext(r.Context())
		if !ok || gotID != id {
			t.Fatalf("expected user id in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
}

func TestAuthRouteGroupPublicVsProtected(t *testing.T) {
	id := "user-1"
	h := newTestHandler(&mockAuthService{signupID: &id, getByIDUser: user.User{ID: &id, Name: "Ada", Email: "ada@example.com"}}, &mockOAuthServer{userID: id})

	r := chi.NewRouter()
	r.Route("/auth", func(rr chi.Router) {
		rr.Post("/signup", h.Signup)
		rr.Group(func(gr chi.Router) {
			gr.Use(RequireBearer(&mockOAuthServer{validErr: errors.New("missing token")}, nil))
			gr.Get("/me", h.Me)
		})
	})

	signupReq := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader([]byte(`{"name":"Ada","email":"ada@example.com","password":"Password1"}`)))
	signupRes := httptest.NewRecorder()
	r.ServeHTTP(signupRes, signupReq)
	if signupRes.Code != http.StatusCreated {
		t.Fatalf("expected signup 201, got %d", signupRes.Code)
	}

	meReq := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	meRes := httptest.NewRecorder()
	r.ServeHTTP(meRes, meReq)
	if meRes.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected route 401 without token, got %d", meRes.Code)
	}
}
