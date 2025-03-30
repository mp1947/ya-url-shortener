package service

import (
	"fmt"
	"strconv"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
)

type Service interface {
	ShortenURL(cfg config.Config, url string) (string, error)
	GetOriginalURL(shortURLID string) (string, error)
	ShortenURLBatch(cfg config.Config, batchData []dto.BatchShortenRequest) ([]dto.BatchShortenResponse, error)
}

type ShortenService struct {
	Storage repository.Repository
	EP      eventlog.EventProcessor
}

func (s *ShortenService) ShortenURL(cfg config.Config, url string) (string, error) {
	ShortURLID := usecase.GenerateIDFromURL(url)
	isSaved, err := s.Storage.Save(ShortURLID, url)

	if err != nil {
		return "", err
	}

	if isSaved {
		s.EP.IncrementUUID()
		event := eventlog.Event{
			UUID:        strconv.Itoa(s.EP.CurrentUUID),
			ShortURL:    ShortURLID,
			OriginalURL: url,
		}
		s.EP.WriteEvent(&event)
	}

	return fmt.Sprintf("%s/%s", *cfg.BaseURL, ShortURLID), nil
}

func (s *ShortenService) ShortenURLBatch(
	cfg config.Config,
	batchData []dto.BatchShortenRequest,
) ([]dto.BatchShortenResponse, error) {

	urls := make([]entity.URL, len(batchData))
	result := make([]dto.BatchShortenResponse, len(batchData))

	for i, v := range batchData {
		shortURLID := usecase.GenerateIDFromURL(v.OriginalURL)
		urls[i] = entity.URL{
			ShortURLID:    shortURLID,
			OriginalURL:   v.OriginalURL,
			CorrelationID: v.CorrelationID,
		}
		result[i] = dto.BatchShortenResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", *cfg.BaseURL, shortURLID),
		}
	}
	_, err := s.Storage.SaveBatch(urls)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *ShortenService) GetOriginalURL(shortURLID string) (string, error) {
	data, err := s.Storage.Get(shortURLID)
	if err != nil {
		return "", err
	}
	return data, nil
}
