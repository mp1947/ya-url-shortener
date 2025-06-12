package dto

// ShortenRequest represents a request payload for shortening a URL.
type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

// BatchShortenRequest represents a request to shorten a URL with a correlation ID.
type BatchShortenRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortenResponse represents the response containing the shortened URL result.
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchShortenResponse represents the response for a batch URL shortening request,
// containing the correlation ID and the resulting short URL.
type BatchShortenResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenURLsByUserID represents a mapping between a shortened URL and its original URL.
type ShortenURLsByUserID struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
