package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/user/create", nil).Methods(http.MethodGet)
}
