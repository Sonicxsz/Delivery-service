package fs

type FS struct {
	Image *Image
}

func New(config *Config) *FS {
	return &FS{
		Image: NewImage(config.Image),
	}
}
