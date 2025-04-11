package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"go.uber.org/zap"
)

type Repository interface {
	Init(ctx context.Context, cfg config.Config, l *zap.Logger) error
	Save(ctx context.Context, shortURLID, originalURL string) error
	SaveBatch(ctx context.Context, urls []entity.URLWithCorrelation) (bool, error)
	Get(ctx context.Context, shortURL string) (string, error)
	GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserURL, error)
	GetType() string
}

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
