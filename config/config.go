// Package config provides functionality for loading and managing application configuration
// from YAML files, environment variables, and command-line flags.
package config

import (
	"flag"
	"log"
	"net"

	"github.com/spf13/viper"
)

const (
	defaultKeysAreNotFoundErr = "error getting defaults from config"
	tlsSettingsUndefinedErr   = "cert file path or config file path were not defined in values.yaml config file"
	defaultServerAddress      = ":8080"
	defaultBaseURL            = "http://localhost:8080"
	defaultFileStoragePath    = "./output.out"
	defaultCrtFilePath        = "./keys/cert.crt"
	defaultKeyFilePath        = "./keys/key.pem"
)

// Config holds the configuration settings for the application, including
// the server listen address, base URL, file storage path, and database DSN.
// All fields are pointers to strings, allowing for optional configuration values.
type Config struct {
	ServerAddress    *string `mapstructure:"SERVER_ADDRESS"`
	BaseURL          *string `mapstructure:"BASE_URL"`
	FileStoragePath  *string `mapstructure:"FILE_STORAGE_PATH"`
	DatabaseDSN      *string `mapstructure:"DATABASE_DSN"`
	TrustedSubnetRaw *string `mapstructure:"TRUSTED_SUBNET"`
	TrustedSubnet    *net.IPNet
	ConfigFilePath   *string
	ShouldUseTLS     *bool `mapstructure:"ENABLE_HTTPS"`
	TLSConfig        *TLS
}

// TLS holds the tls configuration consists of crt and key files path
type TLS struct {
	CrtFilePath string `json:"crt_file"`
	KeyFilePath string `json:"key_file"`
}

// InitConfig initializes the Config struct by loading configuration values from a YAML file,
// command-line flags, and environment variables. It first attempts to read configuration
// values from a "values.yaml" file located in the "./config" or "../../config" directories.
// If any required configuration keys are missing, the function logs a fatal error.
// The function then sets up command-line flags for server address, base URL, file storage path,
// and database DSN, using the loaded configuration values as defaults. After parsing the flags,
// it checks for corresponding environment variables and, if set, overrides the configuration
// values with those from the environment.
func InitConfig() *Config {
	cfg := &Config{}

	cfg.ServerAddress = new(string)
	cfg.BaseURL = new(string)
	cfg.FileStoragePath = new(string)
	cfg.DatabaseDSN = new(string)
	cfg.ConfigFilePath = new(string)
	cfg.ShouldUseTLS = new(bool)
	cfg.TrustedSubnetRaw = new(string)
	cfg.TrustedSubnet = new(net.IPNet)

	flagServerAddress := flag.String("a", "", "listen address, example: -a :8080, default :8080")
	flagBaseURL := flag.String("b", "", "base url, example: -b http://localhost:8080, default: http://localhost:8080")
	flagConfigFile := flag.String("c", "", "config file, example: -c /path/to/config.json")
	flagFileStoragePath := flag.String("f", "", "storage path (inmemory mode), example: -f ./path/to/storage.txt")
	flagDatabaseDSN := flag.String("d", "", "database dsn, example: -d postgres://app:pass@localhost:5432/app?pool_max_conns=10&pool_max_conn_lifetime=1h30m")
	flagShouldUseTLS := flag.Bool("s", false, "if provided, enables https, example: -s")
	flagTrustedSubnet := flag.String("t", "", "trusted subnet: 192.168.1.1./24")
	flag.Parse()

	v := viper.New()
	v.SetConfigType("json")
	v.SetDefault("SERVER_ADDRESS", defaultServerAddress)
	v.SetDefault("BASE_URL", defaultBaseURL)
	v.SetDefault("FILE_STORAGE_PATH", defaultFileStoragePath)
	v.SetDefault("ENABLE_HTTPS", false)
	v.SetDefault("TRUSTED_SUBNET", "")

	if *flagConfigFile != "" {
		v.SetConfigFile(*flagConfigFile)
		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
	}

	v.AutomaticEnv()

	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	if *flagServerAddress != "" {
		cfg.ServerAddress = flagServerAddress
	}
	if *flagBaseURL != "" {
		cfg.BaseURL = flagBaseURL
	}
	if *flagFileStoragePath != "" {
		cfg.FileStoragePath = flagFileStoragePath
	}
	if *flagDatabaseDSN != "" {
		cfg.DatabaseDSN = flagDatabaseDSN
	}
	if *flagShouldUseTLS {
		cfg.ShouldUseTLS = flagShouldUseTLS
	}

	if *cfg.TrustedSubnetRaw != "" {
		_, ipRange, err := net.ParseCIDR(*cfg.TrustedSubnetRaw)
		if err != nil {
			log.Fatalf("not a valid ip range in a TRUSTED_SUBNET variable: %s", *cfg.TrustedSubnetRaw)
		}
		cfg.TrustedSubnet = ipRange
	}

	if *flagTrustedSubnet != "" {
		_, ipRange, err := net.ParseCIDR(*flagTrustedSubnet)

		if err != nil {
			log.Fatalf("not a valid ip range: %s", *flagTrustedSubnet)
		}
		cfg.TrustedSubnet = ipRange
	}

	if *cfg.ShouldUseTLS {
		crtFilePath := viper.GetString("tls_crt_file")
		keyFilePath := viper.GetString("tls_key_file")
		tlsConfig := &TLS{
			CrtFilePath: crtFilePath,
			KeyFilePath: keyFilePath,
		}
		if crtFilePath == "" || keyFilePath == "" {
			log.Printf("tls_crt_file  or tls_key_file not found in config file, setting default values")
			tlsConfig.CrtFilePath = defaultCrtFilePath
			tlsConfig.KeyFilePath = defaultKeyFilePath
		}
		cfg.TLSConfig = tlsConfig
	}

	return cfg
}
