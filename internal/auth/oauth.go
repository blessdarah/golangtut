package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	oauth2 "github.com/go-oauth2/oauth2/v4"
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

func NewOAuthServer(clientID, clientSecret string, accessTokenTTLMinutes, refreshTokenTTLHours int, svc *Service) (OAuthServer, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("missing oauth client configuration")
	}
	if accessTokenTTLMinutes <= 0 {
		return nil, errors.New("invalid oauth access token ttl")
	}
	if refreshTokenTTLHours <= 0 {
		return nil, errors.New("invalid oauth refresh token ttl")
	}

	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    time.Duration(accessTokenTTLMinutes) * time.Minute,
		RefreshTokenExp:   time.Duration(refreshTokenTTLHours) * time.Hour,
		IsGenerateRefresh: true,
	})

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
	srv.SetAllowedGrantType(oauth2.PasswordCredentials)
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
