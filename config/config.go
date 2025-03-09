package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultKeysAreNotFoundErr = "defaults.listen_addr or defaults.base_url is empty"
)

type Config struct {
	ListenAddr *string
	BaseURL    *string
}

func (cfg *Config) ParseFlags() {
	viper.SetConfigName("values")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	defaultListenAddr := viper.GetString("defaults.listen_addr")
	defaultBaseURL := viper.GetString("defaults.base_url")

	if defaultListenAddr == "" || defaultBaseURL == "" {
		log.Fatalf(
			"error reading settings from config: %s",
			errors.New(defaultKeysAreNotFoundErr),
		)
	}

	cfg.ListenAddr = flag.String("a", defaultListenAddr, "-a :8080")
	cfg.BaseURL = flag.String("b", defaultBaseURL, "-b http://localhost:8080")
	flag.Parse()

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ListenAddr = &addr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = &baseURL
	}
}
