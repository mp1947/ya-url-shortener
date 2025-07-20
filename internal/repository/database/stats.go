package database

import (
	"context"

	"github.com/mp1947/ya-url-shortener/internal/dto"
)

// GetInternalStats retrieves internal statistics from the database, including the total number of URLs and users.
// It executes a query to fetch these statistics and returns them as a dto.InternalStatsResp object.
// If an error occurs during the query or scanning process, it returns the error.
func (d *Database) GetInternalStats(ctx context.Context) (*dto.InternalStatsResp, error) {

	data := d.conn.QueryRow(ctx, getInternalStatsQuery)

	result := dto.InternalStatsResp{}

	err := data.Scan(&result.URLs, &result.Users)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
