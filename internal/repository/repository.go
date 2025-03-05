package repository

type Repository interface {
	Save(shortURL, originalURL string)
	Get(shortURL string) string
}
