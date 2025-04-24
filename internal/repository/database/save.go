package database

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

func (d *Database) Save(
	ctx context.Context,
	shortURLID, originalURL string,
	userID string,
) error {

	args := pgx.NamedArgs{
		"shortURL":    shortURLID,
		"originalURL": originalURL,
		"userID":      userID,
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

func (d *Database) SaveBatch(
	ctx context.Context,
	urls []model.URLWithCorrelation,
	userID string,
) (bool, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return false, err
	}

	for _, v := range urls {
		args := pgx.NamedArgs{
			"shortURL":    v.ShortURLID,
			"originalURL": v.OriginalURL,
			"userID":      userID,
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
