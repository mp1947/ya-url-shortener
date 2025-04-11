package database

import (
	"context"
	"embed"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Database struct {
	conn        *pgxpool.Pool
	cfg         config.Config
	StorageType string
}

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

func (d *Database) SaveBatch(ctx context.Context, urls []entity.URLWithCorrelation) (bool, error) {
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
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return false, rbErr
			}
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

func (d *Database) GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserURL, error) {

	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := d.conn.Query(ctx, getURLsByUserID, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	UserURL := make([]entity.UserURL, len(rows.RawValues()))

	for rows.Next() {
		var originalURL, shortURL string

		err := rows.Scan(&originalURL, &shortURL)

		if err != nil {
			return nil, err
		}
		UserURL = append(UserURL, entity.UserURL{
			ShortURLID:  shortURL,
			OriginalURL: originalURL,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserURL, nil
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
