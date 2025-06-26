package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	embed "github.com/mp1947/ya-url-shortener"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// Init initializes the Database connection using the provided configuration and logger.
// It parses the database DSN, establishes a new connection pool, pings the database to ensure connectivity,
// and applies any pending migrations using Goose. If a previous connection exists, it is closed before
// establishing a new one. The function sets the storage type to "database" upon successful initialization.
// Returns an error if any step fails.
func (d *Database) Init(
	ctx context.Context,
	cfg config.Config,
	l *zap.Logger,
) error {
	var err error
	d.cfg = cfg
	pgConfig, err := pgxpool.ParseConfig(*d.cfg.DatabaseDSN)

	if err != nil {
		return err
	}

	if d.conn != nil {
		d.conn.Close()
	}

	d.conn, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return err
	}

	if err := d.conn.Ping(ctx); err != nil {
		return err
	}

	goose.SetBaseFS(embed.EmbedMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(d.conn)
	defer func() {
		err := db.Close()
		if err != nil {
			l.Error("error closing database connection", zap.Error(err))
		}
	}()

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	d.StorageType = "database"
	return nil
}
