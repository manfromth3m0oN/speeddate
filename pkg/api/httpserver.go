package api

import (
	"context"
	"crypto/rsa"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gorilla/mux"
	"github.com/manfromth3m0oN/speeddate/cmd/config"
)

// HTTPService holds the state for the HTTP Server
type HTTPService struct {
	DB      *goqu.Database
	PrivKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
	JWTExpr time.Duration
}

// StartHTTPServer initializes resources and starts a http server as defined by the config
func (h *HTTPService) StartHTTPServer(ctx context.Context, wg *sync.WaitGroup, cfg config.Config) {
	defer wg.Done()
	mux := mux.NewRouter()

	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port,
		Handler: mux,
	}

	go func() {
		slog.Info("starting http server")
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("unexpected http server error", slog.Any("err", err))
		}
	}()

	<-ctx.Done()
	slog.Info("shutting http server down")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("error shutting down server", slog.Any("err", err))
	}

	slog.Info("server shutdown complete")
}
