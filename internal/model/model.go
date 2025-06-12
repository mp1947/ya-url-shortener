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

type BatchDeleteShortURLs struct {
	ShortURLs []string
	UserID    string
}
