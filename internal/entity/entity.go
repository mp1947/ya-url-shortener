package entity

type URLWithCorrelation struct {
	ShortURLID    string
	OriginalURL   string
	CorrelationID string
}

type UserURL struct {
	ShortURLID  string
	OriginalURL string
}

type URL struct {
	ShortURLID  string
	OriginalURL string
	IsDeleted   bool
}

type BatchDeleteShortURLs struct {
	ShortURLs []string
	UserID    string
}
