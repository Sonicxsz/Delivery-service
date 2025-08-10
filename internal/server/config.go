package server

import (
	"arabic/internal/middlewares"
	"arabic/store"
)

type Config struct {
	BindAddr string `toml: "bind_addr"`
	LogLevel string `toml: "log_level"`
	Storage  *store.Config
	JWT      *middlewares.JWTConfig
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Storage:  store.NewConfig(),
		JWT:      middlewares.NewJWTConfig(),
	}
}
