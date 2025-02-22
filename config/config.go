package config

import (
	"flag"
)

// Config is an application config
type Config struct {
	ListenAddr    *string
	BaseResultUrl *string
}

func (c *Config) ParseFlags() {
	c.ListenAddr = flag.String("a", ":8080", "-a :8080")
	c.BaseResultUrl = flag.String("b", "http://localhost:8080", "-b http://localhost:8080")
	flag.Parse()
}
