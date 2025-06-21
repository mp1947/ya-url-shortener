package router_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateRouter(t *testing.T) {
	t.Run("create test router", func(t *testing.T) {
		listenAddr := ":8080"
		baseURL := "http://localhost:8080"
		fileStoragePath := "./test.out"
		cfg := config.Config{
			ListenAddr:      &listenAddr,
			BaseURL:         &baseURL,
			FileStoragePath: &fileStoragePath,
		}
		storage := &inmemory.Memory{}
		l, err := logger.InitLogger()
		assert.NoError(t, err)
		err = storage.Init(context.Background(), cfg, l)
		assert.NoError(t, err)
		service := service.ShortenService{Storage: storage, Logger: l, Cfg: &cfg}
		r := router.CreateRouter(cfg, &service, storage, l)
		assert.IsType(t, &gin.Engine{}, r)
	})
}
