package model

// URLWithCorrelation represents a URL mapping with an associated correlation ID.
// It contains the shortened URL identifier, the original URL, and a correlation ID
// used for tracking or associating requests.
type URLWithCorrelation struct {
	ShortURLID    string
	OriginalURL   string
	CorrelationID string
}

// UserURL represents a mapping between a shortened URL identifier and its original URL.
type UserURL struct {
	ShortURLID  string
	OriginalURL string
}

// URL represents a shortened URL entry with its unique identifier, the original URL,
// and a flag indicating whether the URL has been deleted.
type URL struct {
	ShortURLID  string
	OriginalURL string
	IsDeleted   bool
}

// BatchDeleteShortURLs represents a request to delete a batch of shortened URLs
// associated with a specific user. It contains a slice of short URL identifiers
// and the user ID of the owner.
type BatchDeleteShortURLs struct {
	ShortURLs []string
	UserID    string
}
