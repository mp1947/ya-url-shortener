package repository_test

import (
	"context"
	"testing"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateRepository(t *testing.T) {
	l, err := logger.InitLogger()
	assert.NoError(t, err)

	cfg := config.InitConfig()

	_, err = repository.CreateRepository(l, *cfg, context.Background())
	assert.NoError(t, err)

}
