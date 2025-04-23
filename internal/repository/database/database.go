package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mp1947/ya-url-shortener/config"
)

type Database struct {
	conn        *pgxpool.Pool
	cfg         config.Config
	StorageType string
}

func (d *Database) GetType() string {
	return d.StorageType
}

func (d *Database) Close() {
	d.conn.Close()
}

func (d *Database) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}
