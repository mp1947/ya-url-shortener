package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mp1947/ya-url-shortener/config"
)

type Database struct {
	conn        *pgxpool.Pool
	StorageType string
}

func (d *Database) Init(cfg config.Config) error {
	var err error
	d.conn, err = pgxpool.New(context.TODO(), *cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(context.TODO(), "CREATE TABLE IF NOT EXISTS urls "+
		"(uuid SERIAL PRIMARY KEY, short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) NOT NULL)",
	)
	if err != nil {
		return err
	}
	d.StorageType = "database"
	return nil
}
func (d *Database) Save(shortURL, originalURL string) (bool, error) {
	data, err := d.Get(shortURL)
	if !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	if data != "" {
		return false, nil
	}

	query := `INSERT INTO urls (short_url, original_url) VALUES (@shortURL, @originalURL)`
	args := pgx.NamedArgs{
		"shortURL":    shortURL,
		"originalURL": originalURL,
	}

	_, err = d.conn.Exec(context.TODO(), query, args)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (d *Database) Get(shortURL string) (string, error) {
	query := "SELECT original_url FROM urls where short_url = @shortURL"
	args := pgx.NamedArgs{
		"shortURL": shortURL,
	}
	row := d.conn.QueryRow(context.TODO(), query, args)
	var shortURLFromDB string
	err := row.Scan(&shortURLFromDB)
	if err != nil {
		return "", err
	}
	return shortURLFromDB, nil
}

func (d *Database) GetType() string {
	return d.StorageType
}

func (d *Database) Close() {
	d.conn.Close()
}

func (d *Database) Ping() error {
	return d.conn.Ping(context.TODO())
}
