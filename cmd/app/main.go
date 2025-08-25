package main

import (
	"arabic/internal/server"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "Path to api server config")
}

func main() {
	flag.Parse()
	config := server.NewConfig()

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		println("Cannot get config file, using default values")
	}

	if err := config.Storage.RunMigrations(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	server := server.New(config)

	println("Server starting")

	log.Fatal(server.Start())
}
