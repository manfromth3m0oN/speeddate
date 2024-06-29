package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/manfromth3m0oN/speeddate/cmd/config"
	"github.com/manfromth3m0oN/speeddate/pkg/api"
)

var programLevel = new(slog.LevelVar) // Info by default

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init Config
	cfg, err := config.BuildConfig()
	if err != nil {
		slog.Error("Failed to init config", slog.Any("err", err))
		os.Exit(1)
	}

	// Set the default logger based on the config
	buildLogger(cfg)

	// Start HTTP Server and wait for it to finish
	var wg sync.WaitGroup
	wg.Add(1)
	go api.StartHTTPServer(ctx, &wg, cfg)

	// Create channel for termination signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	cancelFunc()
	wg.Wait()
}

func buildLogger(cfg config.Config) {
	// Set log type based on config
	var handler slog.Handler
	switch cfg.Logging.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	case "term":
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	}
	slog.SetDefault(slog.New(handler))

	// Set default log level based on config
	switch cfg.Logging.Level {
	case "error":
		programLevel.Set(slog.LevelError)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	}
}
