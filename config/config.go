package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultKeysAreNotFoundErr = "error getting defaults from config"
)

// Config holds the configuration settings for the application, including
// the server listen address, base URL, file storage path, and database DSN.
// All fields are pointers to strings, allowing for optional configuration values.
type Config struct {
	ListenAddr      *string
	BaseURL         *string
	FileStoragePath *string
	DatabaseDSN     *string
}

// InitConfig initializes the Config struct by loading configuration values from a YAML file,
// command-line flags, and environment variables. It first attempts to read configuration
// values from a "values.yaml" file located in the "./config" or "../../config" directories.
// If any required configuration keys are missing, the function logs a fatal error.
// The function then sets up command-line flags for server address, base URL, file storage path,
// and database DSN, using the loaded configuration values as defaults. After parsing the flags,
// it checks for corresponding environment variables and, if set, overrides the configuration
// values with those from the environment.
func (cfg *Config) InitConfig() {
	viper.SetConfigName("values")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	defaultListenAddr := viper.GetString("defaults.listen_addr")
	defaultBaseURL := viper.GetString("defaults.base_url")
	defaultFileStoragePath := viper.GetString("defaults.file_storage_path")

	if defaultListenAddr == "" || defaultBaseURL == "" || defaultFileStoragePath == "" {
		log.Fatalf(
			"error reading settings from config: %s",
			errors.New(defaultKeysAreNotFoundErr),
		)
	}

	cfg.ListenAddr = flag.String("a", defaultListenAddr, "-a :8080")
	cfg.BaseURL = flag.String("b", defaultBaseURL, "-b http://localhost:8080")
	cfg.FileStoragePath = flag.String("f", defaultFileStoragePath, "-f ./storage/storage.txt")
	cfg.DatabaseDSN = flag.String("d", "", "-d postgres://app:pass@localhost:5432/app?pool_max_conns=10&pool_max_conn_lifetime=1h30m")
	flag.Parse()

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ListenAddr = &addr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = &baseURL
	}

	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		cfg.FileStoragePath = &fileStoragePath
	}

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		cfg.DatabaseDSN = &databaseDSN
	}
}
