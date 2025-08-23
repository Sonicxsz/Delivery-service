package server

import (
	security "arabic/internal/security/auth"
	"arabic/store"
)

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	LogDir   string `toml:"log_dir"`
	Storage  *store.Config
	JWT      *security.JWTConfig
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Storage:  store.NewConfig(),
		JWT:      security.NewJWTConfig(),
	}
}
