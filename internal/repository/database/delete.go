package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (d *Database) DeleteBatch(
	ctx context.Context,
	shortIDs []string,
	userID string,
) error {

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, v := range shortIDs {
		args := pgx.NamedArgs{
			"shortURL": v,
			"userID":   userID,
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
