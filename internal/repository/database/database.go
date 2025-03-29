package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mp1947/ya-url-shortener/config"
)

type Database struct {
	conn        *pgx.Conn
	StorageType string
}

func (d *Database) Init(cfg config.Config) error {
	var err error
	d.conn, err = pgx.Connect(context.TODO(), *cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	d.StorageType = "database"
	return nil
}
func (d *Database) Save(shortURL, originalURL string) bool {
	return false
}
func (d *Database) Get(shortURL string) string {
	return ""
}

func (d *Database) GetType() string {
	return ""
}

func (d *Database) Close() error {
	return d.conn.Close(context.TODO())
}

func (d *Database) Ping() error {
	return d.conn.Ping(context.TODO())
}
