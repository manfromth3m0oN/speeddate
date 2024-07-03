package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kkyr/fig"
)

var (
	ErrMissingEnv = errors.New("missing env var")
)

type EnvVar string

const (
	Env     EnvVar = "ENV"
	AppName EnvVar = "APP_NAME"
)

// Lookup an environment variable
func (e EnvVar) Lookup() (string, bool) {
	return os.LookupEnv(string(e))
}

// Config defines the structure of the config file
type Config struct {
	Logging struct {
		Level  string `fig:"level" default:"info"`
		Format string `fig:"format" default:"term"`
	}
	HTTPServer struct {
		Host       string        `fig:"host" default:"localhost"`
		Port       string        `fig:"port" default:"8080"`
		JWTExpr    time.Duration `fig:"jwt_expiration"`
		JWTPubKey  string        `fig:"jwt_public_key"`
		JWTPrivKey string        `fig:"jwt_private_key"`
	} `fig:"http_server"`
	Database struct {
		User     string `fig:"user"`
		Password string `fig:"password"`
		Host     string `fig:"host"`
		Port     string `fig:"port"`
	}
}

// BuildConfig reads in the config file and builds it into the struct
func BuildConfig() (Config, error) {
	var config Config

	env, envPresent := Env.Lookup()
	name, namePresent := AppName.Lookup()
	if !envPresent || !namePresent {
		slog.Error("missing env var", slog.Any("envPresent", envPresent), slog.Any("namePresent", namePresent))
		return config, ErrMissingEnv
	}

	configName := fmt.Sprintf("%s_%s.toml", env, name)

	if err := fig.Load(&config, fig.File(configName), fig.Dirs("./conf")); err != nil {
		return Config{}, err
	}
	slog.Debug("config loaded", slog.Any("cfg", config))
	return config, nil
}
