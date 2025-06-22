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

// Save inserts a new short URL mapping into the database, associating the given shortURLID with the originalURL and userID.
// If the originalURL already exists for a user, it returns shrterr.ErrOriginalURLAlreadyExists.
// Returns an error if the operation fails for other reasons.
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
		if pgErr.Code == pgerrcode.UniqueViolation {
			return shrterr.ErrOriginalURLAlreadyExists
		}
	} else if err != nil {
		return err
	}

	return nil
}

// SaveBatch saves a batch of shortened URLs to the database within a single transaction.
// It takes a context, a slice of model.URLWithCorrelation containing the URLs to save,
// and the userID associated with the URLs. If any insert fails, the transaction is rolled back.
// Returns true if all URLs are saved successfully, otherwise returns false and an error.
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
		_, ExecErr := tx.Exec(ctx, insertShortURLQuery, args)
		if ExecErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return false, rbErr
			}
			return false, ExecErr
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}
