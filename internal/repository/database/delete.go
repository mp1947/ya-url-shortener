package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mp1947/ya-url-shortener/internal/entity"
)

func (d *Database) DeleteBatch(
	ctx context.Context,
	shortURLs entity.BatchDeleteShortURLs,
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
