package fs

type Config struct {
	Image *ImageConfig `toml:"image"`
	Path  string       `toml:"static_path"`
}

func NewFSConfig() *Config {
	return &Config{
		Image: NewImageConfig(),
	}
}
