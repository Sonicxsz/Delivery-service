package main

import (
	"arabic/internal/server"
	"arabic/store"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "Path to api server config")
}

func runMigrations(config *store.Config) error {

	m, err := migrate.New(
		config.DbMigrationsPath,
		config.DbMigrationsUrl,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	config := server.NewConfig()

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		println("Cannot get config file, using default values")
	}

	if err := runMigrations(config.Storage); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	server := server.New(config)

	println("Server starting")
	log.Fatal(server.Start())
}
