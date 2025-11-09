package config

// Config holds simple configuration (placeholder).
type Config struct {
	Addr string
}

func Default() *Config { return &Config{Addr: ":8443"} }
