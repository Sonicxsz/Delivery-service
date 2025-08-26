package main

import (
	"arabic/internal/server"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
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

	api := server.New(config)

	println("Server starting")

	log.Fatal(api.Start())
}
