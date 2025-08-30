package store

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	DbConnString     string `toml:"db_conn_string"`
	DbMigrationsUrl  string `toml:"db_migrations_url"`
	DbMigrationsPath string `toml:"db_migrations_path"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) RunMigrations() error {
	m, err := migrate.New(
		c.DbMigrationsPath,
		c.DbMigrationsUrl,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
