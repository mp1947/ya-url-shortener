package database

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
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
	_, err = d.conn.Exec(context.TODO(), createTableQuery)
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(context.TODO(), createIndexQuery)
	if err != nil {
		return err
	}

	d.StorageType = "database"
	return nil
}
func (d *Database) Save(shortURLID, originalURL string) error {

	args := pgx.NamedArgs{
		"shortURL":    shortURLID,
		"originalURL": originalURL,
	}

	_, err := d.conn.Exec(context.TODO(), insertShortURLQuery, args)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return shrterr.ErrOriginalURLAlreadyExists
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (d *Database) SaveBatch(urls []entity.URL) (bool, error) {
	tx, err := d.conn.Begin(context.TODO())
	if err != nil {
		return false, err
	}

	for _, v := range urls {
		args := pgx.NamedArgs{
			"shortURL":    v.ShortURLID,
			"originalURL": v.OriginalURL,
		}
		_, err := tx.Exec(context.TODO(), insertShortURLQuery, args)
		if err != nil {
			tx.Rollback(context.TODO())
			return false, err
		}
	}

	err = tx.Commit(context.TODO())

	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Database) Get(shortURL string) (string, error) {
	args := pgx.NamedArgs{
		"shortURL": shortURL,
	}
	row := d.conn.QueryRow(context.TODO(), getOriginalURLByShortIDQuery, args)
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
