package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/manfromth3m0oN/speeddate/pkg/user"
)

// CreateUser handles an HTTP request to create a new user
func (h *HTTPService) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		newUser := user.NewRandomUser()

		err := newUser.Insert(ctx, h.DB)
		if err != nil {
			slog.Error("failed to insert new user", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(newUser); err != nil {
			slog.Error("failed to encode new user", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
