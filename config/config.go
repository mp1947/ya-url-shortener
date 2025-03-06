package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddr *string
	BaseURL    *string
}

func (cfg *Config) ParseFlags() {
	cfg.ListenAddr = flag.String("a", ":8080", "-a :8080")
	cfg.BaseURL = flag.String("b", "http://localhost:8080", "-b http://localhost:8080")
	flag.Parse()

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ListenAddr = &addr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = &baseURL
	}
}
