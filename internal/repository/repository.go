package repository

type Repository interface {
	Init()
	Save(shortURL, originalURL string)
	Get(shortURL string) string
}
