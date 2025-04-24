package service

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	ShortenURL(
		ctx context.Context,
		url string,
		userID string,
	) (string, error)
	GetOriginalURL(ctx context.Context, shortURLID string) (model.URL, error)
	ShortenURLBatch(
		ctx context.Context,
		batchData []dto.BatchShortenRequest,
		userID string,
	) ([]dto.BatchShortenResponse, error)
	DeleteURLsBatch(ctx context.Context, shortURLs model.BatchDeleteShortURLs)
	GetUserURLs(
		ctx context.Context,
		userID string,
	) ([]dto.ShortenURLsByUserID, error)
}

type ShortenService struct {
	Storage repository.Repository
	EP      eventlog.EventProcessor
	Cfg     *config.Config
	Logger  *zap.Logger
	CommCh  chan model.BatchDeleteShortURLs
}
