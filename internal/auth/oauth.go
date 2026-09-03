package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

type OAuthServer interface {
	HandleTokenRequest(w http.ResponseWriter, r *http.Request) error
	ValidateBearerToken(r *http.Request) (string, error)
}

type oauthServer struct {
	server *server.Server
}

func NewOAuthServer(clientID, clientSecret string, svc *Service) (OAuthServer, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing oauth client configuration")
	}

	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	err := clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
	})
	if err != nil {
		return nil, err
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (string, error) {
		return svc.ValidateCredentials(username, password)
	})

	return &oauthServer{server: srv}, nil
}

func (s *oauthServer) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.server.HandleTokenRequest(w, r)
}

func (s *oauthServer) ValidateBearerToken(r *http.Request) (string, error) {
	tokenInfo, err := s.server.ValidationBearerToken(r)
	if err != nil {
		return "", err
	}

	return tokenInfo.GetUserID(), nil
}
