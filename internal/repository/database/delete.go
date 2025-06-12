package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

// DeleteBatch deletes a batch of short URLs associated with a specific user from the database.
// It starts a transaction and iterates over the provided short URLs, executing a delete query for each.
// If any deletion fails, the transaction is rolled back and the error is returned.
// On success, the transaction is committed and the number of rows affected by the last delete operation is returned.
// Parameters:
//   - ctx: context for controlling cancellation and deadlines.
//   - shortURLs: a BatchDeleteShortURLs struct containing the user ID and a slice of short URLs to delete.
//
// Returns:
//   - int64: the number of rows affected by the last delete operation.
//   - error: an error if the operation fails, otherwise nil.
func (d *Database) DeleteBatch(
	ctx context.Context,
	shortURLs model.BatchDeleteShortURLs,
) (int64, error) {

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return 0, err
	}

	var ct pgconn.CommandTag

	for _, v := range shortURLs.ShortURLs {
		args := pgx.NamedArgs{
			"shortURL": v,
			"userID":   shortURLs.UserID,
		}
		ct, err = tx.Exec(ctx, deleteURLQuery, args)
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return 0, rbErr
			}
			return 0, err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}
