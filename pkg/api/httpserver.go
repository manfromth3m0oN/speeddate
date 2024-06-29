package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/manfromth3m0oN/speeddate/cmd/config"
)

func StartHTTPServer(ctx context.Context, wg *sync.WaitGroup, cfg config.Config) {
	defer wg.Done()
	mux := mux.NewRouter()

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
