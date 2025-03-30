package repository

import "github.com/mp1947/ya-url-shortener/config"

type Repository interface {
	Init(cfg config.Config) error
	Save(shortURL, originalURL string) (bool, error)
	Get(shortURL string) (string, error)
	Ping() error
	GetType() string
}
