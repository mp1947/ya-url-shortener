package database

const (
	insertShortURLQuery = `
	INSERT INTO urls (short_url, original_url, user_uuid) 
	VALUES (@shortURL, @originalURL, @userID)
	`
	getOriginalURLByShortIDQuery = `SELECT original_url, is_deleted FROM urls where short_url = @shortURL`
	getURLsByUserID              = `SELECT original_url, short_url FROM urls where user_uuid = @userID`
	deleteURLQuery               = `UPDATE urls SET is_deleted = true WHERE short_url = @shortURL AND user_uuid = @userID`
	getInternalStatsQuery        = `SELECT count(*), count(distinct user_uuid) from urls`
)
