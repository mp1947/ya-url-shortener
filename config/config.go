package config

import (
	"flag"
)

// Config is an application config
type Config struct {
	ListenAddr *string
	BasePath   *string
}

func (c *Config) ParseFlags() {
	c.ListenAddr = flag.String("a", ":8080", "-a 127.0.0.1:8080")
	c.BasePath = flag.String("b", "/", "-b /path")
	flag.Parse()
}
