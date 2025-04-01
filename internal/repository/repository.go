package repository

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
)

type Repository interface {
	Init(cfg config.Config, ctx context.Context) error
	Save(shortURLID, originalURL string) error
	SaveBatch(urls []entity.URL) (bool, error)
	Get(shortURL string) (string, error)
	Ping() error
	GetType() string
	RestoreFromFile() (int, error)
}
