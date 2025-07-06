package config_test

import (
	"testing"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	t.Run("test_init_config", func(t *testing.T) {
		cfg := config.InitConfig()

		assert.NotNil(t, cfg.BaseURL)
		assert.NotNil(t, cfg.DatabaseDSN)
		assert.NotNil(t, cfg.FileStoragePath)
		assert.NotNil(t, cfg.ServerAddress)
	})

}
