package database

const (
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS urls (
    uuid SERIAL PRIMARY KEY,
    short_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255) NOT NULL
	)
	`
	createIndexQuery = `
	CREATE UNIQUE INDEX IF NOT EXISTS 
		original_url ON urls (original_url)
	`
	insertShortURLQuery = `
	INSERT INTO urls (short_url, original_url) 
	VALUES (@shortURL, @originalURL)
	`
	getOriginalURLByShortIDQuery = `SELECT original_url FROM urls where short_url = @shortURL`
)
