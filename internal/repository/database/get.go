package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mp1947/ya-url-shortener/internal/model"
)

func (d *Database) Get(ctx context.Context, shortURL string) (model.URL, error) {
	args := pgx.NamedArgs{
		"shortURL": shortURL,
	}
	row := d.conn.QueryRow(ctx, getOriginalURLByShortIDQuery, args)
	var originalURLFromDB string
	var isDeleted bool
	err := row.Scan(&originalURLFromDB, &isDeleted)
	if err != nil {
		return model.URL{}, err
	}
	return model.URL{
		ShortURLID:  shortURL,
		OriginalURL: originalURLFromDB,
		IsDeleted:   isDeleted,
	}, nil
}

// GetURLsByUserID retrieves all URLs associated with the specified user ID from the database.
// It returns a slice of model.UserURL containing the original and shortened URLs for the user.
// If an error occurs during the query or scanning process, it returns the error.
func (d *Database) GetURLsByUserID(ctx context.Context, userID string) ([]model.UserURL, error) {

	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := d.conn.Query(ctx, getURLsByUserID, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	UserURL := make([]model.UserURL, len(rows.RawValues()))

	for rows.Next() {
		var originalURL, shortURL string

		err := rows.Scan(&originalURL, &shortURL)

		if err != nil {
			return nil, err
		}
		UserURL = append(UserURL, model.UserURL{
			ShortURLID:  shortURL,
			OriginalURL: originalURL,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UserURL, nil
}
