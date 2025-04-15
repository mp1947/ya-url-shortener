package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mp1947/ya-url-shortener/internal/entity"
)

func (d *Database) Get(ctx context.Context, shortURL string) (entity.URL, error) {
	args := pgx.NamedArgs{
		"shortURL": shortURL,
	}
	row := d.conn.QueryRow(ctx, getOriginalURLByShortIDQuery, args)
	var originalURLFromDB string
	var isDeleted bool
	err := row.Scan(&originalURLFromDB, &isDeleted)
	if err != nil {
		return entity.URL{}, err
	}
	return entity.URL{
		ShortURLID:  shortURL,
		OriginalURL: originalURLFromDB,
		IsDeleted:   isDeleted,
	}, nil
}

func (d *Database) GetURLsByUserID(ctx context.Context, userID string) ([]entity.UserURL, error) {

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
