// Package service provides core business logic for URL shortening operations,
// including URL creation, retrieval, batch processing, deletion, and user-specific queries.
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

// Service defines the interface for URL shortening service operations.
// It provides methods for shortening URLs (individually and in batch),
// retrieving the original URL by its shortened ID, deleting batches of URLs,
// and fetching all shortened URLs associated with a specific user.
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

// ShortenService provides methods for URL shortening operations.
// It manages storage, event processing, configuration, logging, and batch deletion communication.
// Fields:
//   - Storage: Interface to the URL repository for storing and retrieving shortened URLs.
//   - EP: Event processor for handling service events.
//   - Cfg: Service configuration settings.
//   - Logger: Structured logger for service logging.
//   - CommCh: Channel for batch deletion of short URLs.
type ShortenService struct {
	Cfg     *config.Config
	Logger  *zap.Logger
	CommCh  chan model.BatchDeleteShortURLs
	Storage repository.Repository
	EP      eventlog.EventProcessor
}
