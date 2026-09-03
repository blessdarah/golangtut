package auth

import (
	"blessdarah/tuts/internal/lib"
	"context"
	"log/slog"
	"net/http"
)

type authContextKey string

const userIDContextKey authContextKey = "auth.user_id"

func userIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDContextKey)
	userID, ok := v.(string)
	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	return userIDFromContext(ctx)
}

func withUserID(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), userIDContextKey, userID)
	return r.WithContext(ctx)
}

func RequireBearer(oauth OAuthServer, logger *slog.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := oauth.ValidateBearerToken(r)
			if err != nil {
				if logger != nil {
					logger.Error("validate bearer token", "error", err)
				}
				lib.WriteProblem(w, r, lib.ProblemDetails{
					Type:   lib.ProblemTypeValidationError,
					Title:  "Unauthorized",
					Status: http.StatusUnauthorized,
					Detail: "invalid or missing access token",
				})
				return
			}

			next.ServeHTTP(w, withUserID(r, userID))
		})
	}

}
