package repository

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	"go.uber.org/zap"
)

type Repository interface {
	Init(cfg config.Config, ctx context.Context) error
	Save(ctx context.Context, shortURLID, originalURL string) error
	SaveBatch(ctx context.Context, urls []entity.URL) (bool, error)
	Get(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) error
	GetType() string
	RestoreFromFile(l *zap.Logger) (int, error)
}
