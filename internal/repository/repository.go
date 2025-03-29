package repository

import "github.com/mp1947/ya-url-shortener/config"

type Repository interface {
	Init(cfg config.Config) error
	Save(shortURL, originalURL string) bool
	Get(shortURL string) string
	Ping() error
	GetType() string
}
