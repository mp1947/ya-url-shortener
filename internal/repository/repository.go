package repository

type Repository interface {
	Init()
	Save(shortURL, originalURL string) bool
	Get(shortURL string) string
	GetType() string
}
