package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type ContextKey int

const UserIDKey ContextKey = iota

// AuthMiddleware makes sure a user is logged in with a token before interacting with certian endpoints
// It also provides the UserIDKey context value, identifying the user
func (h *HTTPService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		slog.Info("got auth header", "header", authHeader[7:])

		token, err := jwt.Parse([]byte(authHeader[7:]), jwt.WithKey(jwa.RS256, h.PubKey))
		if err != nil {
			slog.Error("failed to parse token", "err", err)
			http.Error(w, "failed to parse token", http.StatusInternalServerError)
			return
		}

		if token.Expiration().Before(time.Now()) {
			http.Error(w, "token expired", http.StatusInternalServerError)
			return
		}

		userID, err := strconv.Atoi(token.Audience()[0])
		if err != nil {
			http.Error(w, "user id is not an int", http.StatusBadRequest)
			return
		}

		nextReq := r.WithContext(context.WithValue(ctx, UserIDKey, userID))

		next.ServeHTTP(w, nextReq)
	})
}
