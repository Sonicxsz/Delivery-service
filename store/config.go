package store

type Config struct {
	DbConnString     string `toml:"db_conn_string"`
	DbMigrationsUrl  string `toml:"db_migrations_url"`
	DbMigrationsPath string `toml:"db_migrations_path"`
}

func NewConfig() *Config {
	return &Config{}
}
