package server

import (
	"arabic/internal/store"
	"arabic/pkg/fs"
	"arabic/pkg/security/auth"
)

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	LogDir   string `toml:"log_dir"`
	Storage  *store.Config
	JWT      *security.JWTConfig
	FS       *fs.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Storage:  store.NewConfig(),
		JWT:      security.NewJWTConfig(),
		FS:       fs.NewFSConfig(),
	}
}
