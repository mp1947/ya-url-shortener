package repository

import (
	"context"

	"github.com/mp1947/ya-url-shortener/config"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"go.uber.org/zap"
)

// Repository defines the interface for URL storage and retrieval operations.
// It abstracts the underlying data storage mechanism and provides methods for
// initializing the repository, saving single or multiple URLs, deleting URLs in batch,
// retrieving a URL by its short identifier, fetching all URLs associated with a user,
// and obtaining the repository type.
type Repository interface {
	Init(ctx context.Context, cfg config.Config, l *zap.Logger) error
	Save(ctx context.Context, shortURLID, originalURL string, userID string) error
	SaveBatch(ctx context.Context, urls []model.URLWithCorrelation, userID string) (bool, error)
	DeleteBatch(ctx context.Context, shortURLs model.BatchDeleteShortURLs) (int64, error)
	Get(ctx context.Context, shortURL string) (model.URL, error)
	GetURLsByUserID(ctx context.Context, userID string) ([]model.UserURL, error)
	GetType() string
}

// CreateRepository initializes and returns a storage repository based on the provided configuration.
// It selects either an in-memory or database-backed storage implementation depending on the presence
// of a database DSN in the configuration. For in-memory storage, it attempts to restore records from
// a file if a file storage path is specified. Logs are emitted for key initialization steps.
// Returns the initialized Repository or an error if initialization fails.
func CreateRepository(
	l *zap.Logger,
	cfg config.Config,
	ctx context.Context,
) (Repository, error) {
	l.Info("initializing storage backend")

	var storageType string

	if *cfg.DatabaseDSN != "" {
		storageType = "database"
	} else {
		storageType = "inmemory"
	}

	switch storageType {
	case "inmemory":
		l.Info("creating inmemory storage backend")
		m := &inmemory.Memory{}
		if err := m.Init(ctx, cfg, l); err != nil {
			return nil, err
		}
		l.Info(
			"restoring records from file storage",
			zap.String("file_storage_path", *cfg.FileStoragePath),
		)
		numRecordsRestored, err := m.RestoreFromFile(l)
		if err != nil {
			l.Fatal("error loading data from file", zap.Error(err))
		}
		l.Info("records restored", zap.Int("count", numRecordsRestored))
		return m, nil
	case "database":
		l.Info("creating database storage backend")
		db := &database.Database{}
		if err := db.Init(ctx, cfg, l); err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, shrterr.ErrUnableToDetermineStorageType
}
