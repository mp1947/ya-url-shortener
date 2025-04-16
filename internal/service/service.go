package service

import (
	"context"
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	ShortenURL(
		ctx context.Context,
		cfg config.Config,
		url string,
		userID string,
	) (string, error)
	GetOriginalURL(ctx context.Context, shortURLID string) (entity.URL, error)
	ShortenURLBatch(
		ctx context.Context,
		cfg config.Config,
		batchData []dto.BatchShortenRequest,
		userID string,
	) ([]dto.BatchShortenResponse, error)
	DeleteURLsBatch(ctx context.Context, shortURLs entity.BatchDeleteShortURLs) error
	GetUserURLs(
		ctx context.Context,
		cfg config.Config,
		userID string,
	) ([]dto.ShortenURLsByUserID, error)
}

type ShortenService struct {
	Storage repository.Repository
	EP      eventlog.EventProcessor
	Logger  *zap.Logger
	CommCh  chan entity.BatchDeleteShortURLs
}

func generateShortURL(baseURL, shortURLID string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortURLID)
}
