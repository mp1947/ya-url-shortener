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
