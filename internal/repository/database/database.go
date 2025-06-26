// Package database provides a PostgreSQL-backed storage layer with connection pooling and configuration management.
package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mp1947/ya-url-shortener/config"
)

// Database represents a storage layer backed by a PostgreSQL connection pool.
// It holds the database connection pool, configuration settings, and the type of storage used.
type Database struct {
	conn        *pgxpool.Pool
	cfg         config.Config
	StorageType string
}

// GetType returns the storage type used by the Database as a string.
func (d *Database) GetType() string {
	return d.StorageType
}

// Close closes the database connection held by the Database instance.
// It should be called when the Database is no longer needed to release resources.
func (d *Database) Close() {
	d.conn.Close()
}

// Ping checks the connection to the database by sending a ping request.
// It returns an error if the database is unreachable or the ping fails.
func (d *Database) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}
