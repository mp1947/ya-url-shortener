package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mp1947/ya-url-shortener/internal/entity"
)

func (d *Database) DeleteBatch(
	ctx context.Context,
	shortURLs entity.BatchDeleteShortURLs,
) error {

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, v := range shortURLs.ShortURLs {
		args := pgx.NamedArgs{
			"shortURL": v,
			"userID":   shortURLs.UserID,
		}
		_, err := tx.Exec(ctx, deleteURLQuery, args)
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return rbErr
			}
			return err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}
