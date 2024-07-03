package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/manfromth3m0oN/speeddate/cmd/config"
	"github.com/manfromth3m0oN/speeddate/pkg/api"
	"github.com/manfromth3m0oN/speeddate/pkg/db"
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

	// Connect to the DB
	db, err := db.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to the database", "err", err)
		os.Exit(1)
	}

	// Pull in the keys for the JWTs
	pubK, privK, err := ReadRSAKey(cfg)
	if err != nil {
		slog.Error("failed to read keys", "err", err)
	}

	httpService := api.HTTPService{
		DB:      db,
		PrivKey: privK,
		PubKey:  pubK,
		JWTExpr: cfg.HTTPServer.JWTExpr,
	}

	// Start HTTP Server and wait for it to finish
	var wg sync.WaitGroup
	wg.Add(1)
	go httpService.StartHTTPServer(ctx, &wg, cfg)

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

// ReadRSAKey reads the RSA keys
func ReadRSAKey(cfg config.Config) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	pubFile, err := os.ReadFile(cfg.HTTPServer.JWTPubKey)
	if err != nil {
		return nil, nil, err
	}
	privFile, err := os.ReadFile(cfg.HTTPServer.JWTPrivKey)
	if err != nil {
		return nil, nil, err
	}

	pubBlock, _ := pem.Decode(pubFile)
	privBlock, _ := pem.Decode(privFile)

	pub, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	priv, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return pub.(*rsa.PublicKey), priv.(*rsa.PrivateKey), nil
}
