package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers endpoints with a mux router
func (h *HTTPService) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/create", h.CreateUser()).Methods(http.MethodPost)
	// Ensure login requests come over https so request body is encrypted
	router.HandleFunc("/login", h.LoginUser()).Methods(http.MethodPost) //.Schemes("https")

	router.Handle("/discover", h.AuthMiddleware(h.DiscoverUsers())).Methods(http.MethodGet)
	router.Handle("/swipe", h.AuthMiddleware(h.SwipeOnUser())).Methods(http.MethodPost)
}
