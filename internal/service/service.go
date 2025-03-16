package service

import (
	"fmt"
	"strconv"

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
	EP      EventProcessor
}

func (s *ShortenService) ShortenURL(cfg config.Config, url string) string {
	ShortURLID := usecase.GenerateIDFromURL(url)
	isSaved := s.Storage.Save(ShortURLID, url)

	if isSaved {
		s.EP.incrementUUID()
		event := Event{
			UUID:        strconv.Itoa(s.EP.currentUUID),
			ShortURL:    ShortURLID,
			OriginalURL: url,
		}
		s.EP.writeEvent(&event)
	}

	return fmt.Sprintf("%s/%s", *cfg.BaseURL, ShortURLID)
}

func (s *ShortenService) GetOriginalURL(shortURLID string) string {
	return s.Storage.Get(shortURLID)
}
