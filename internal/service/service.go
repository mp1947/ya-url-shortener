package service

import (
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
)

type Service interface {
	ShortenURL(cfg config.Config, url string) string
	GetOriginalURL(shortURLID string) string
}

type ShortenService struct {
	Storage repository.Repository
}

func (s ShortenService) ShortenURL(cfg config.Config, url string) string {
	ShortURLID := usecase.GenerateIDFromURL(url)
	s.Storage.Save(ShortURLID, url)
	return fmt.Sprintf("%s/%s", *cfg.BaseURL, ShortURLID)
}

func (s ShortenService) GetOriginalURL(shortURLID string) string {
	return s.Storage.Get(shortURLID)
}
