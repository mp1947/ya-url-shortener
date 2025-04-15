package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

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

	d.conn, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return err
	}

	if err := d.conn.Ping(ctx); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(d.conn)
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	d.StorageType = "database"
	return nil
}
