package database

const (
	insertShortURLQuery = `
	INSERT INTO urls (short_url, original_url) 
	VALUES (@shortURL, @originalURL)
	`
	getOriginalURLByShortIDQuery = `SELECT original_url FROM urls where short_url = @shortURL`
	getURLsByUserID              = `SELECT original_url, short_url FROM urls where user_uuid = @userID`
)
