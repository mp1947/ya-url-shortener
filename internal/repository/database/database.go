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
	cfg         config.Config
	StorageType string
}

func (d *Database) Init(cfg config.Config, ctx context.Context) error {
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
	_, err = d.conn.Exec(ctx, createTableQuery)
	if err != nil {
		return err
	}

	_, err = d.conn.Exec(ctx, createIndexQuery)
	if err != nil {
		return err
	}

	d.StorageType = "database"
	return nil
}
func (d *Database) Save(ctx context.Context, shortURLID, originalURL string) error {

	args := pgx.NamedArgs{
		"shortURL":    shortURLID,
		"originalURL": originalURL,
	}

	_, err := d.conn.Exec(ctx, insertShortURLQuery, args)
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

func (d *Database) SaveBatch(ctx context.Context, urls []entity.URL) (bool, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return false, err
	}

	for _, v := range urls {
		args := pgx.NamedArgs{
			"shortURL":    v.ShortURLID,
			"originalURL": v.OriginalURL,
		}
		_, err := tx.Exec(ctx, insertShortURLQuery, args)
		if err != nil {
			tx.Rollback(ctx)
			return false, err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Database) Get(ctx context.Context, shortURL string) (string, error) {
	args := pgx.NamedArgs{
		"shortURL": shortURL,
	}
	row := d.conn.QueryRow(ctx, getOriginalURLByShortIDQuery, args)
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

func (d *Database) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}

func (d *Database) RestoreFromFile() (int, error) {
	return 0, nil
}
