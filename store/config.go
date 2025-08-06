package store

type Config struct {
	DbConnString string `toml:"db_conn_string"`
}

func NewConfig() *Config {
	return &Config{}
}
