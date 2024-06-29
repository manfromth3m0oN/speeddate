package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

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

func (e EnvVar) Lookup() (string, bool) {
	return os.LookupEnv(string(e))
}

type Config struct {
	Logging struct {
		Level  string `fig:"level" default:"info"`
		Format string `fig:"format" default:"term"`
	}
	HTTPServer struct {
		Host string `fig:"host" default:"localhost"`
		Port string `fig:"port" default:"8080"`
	} `fig:"http_server"`
}

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
