package service

import (
	"errors"
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"go.uber.org/zap"
)

type Service interface {
	ShortenURL(cfg config.Config, url string) (string, error)
	GetOriginalURL(shortURLID string) (string, error)
	ShortenURLBatch(cfg config.Config, batchData []dto.BatchShortenRequest) ([]dto.BatchShortenResponse, error)
}

type ShortenService struct {
	Storage repository.Repository
	EP      eventlog.EventProcessor
	Logger  *zap.Logger
}

func (s *ShortenService) ShortenURL(cfg config.Config, url string) (string, error) {
	s.Logger.Info("shortening incoming url", zap.String("original_url", url))

	shortURLID := usecase.GenerateIDFromURL(url)

	s.Logger.Info(
		"short_url id generated for url",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", url),
	)

	err := s.Storage.Save(shortURLID, url)
	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		s.Logger.Info(
			"got an error from repository on save, returning error with short url",
			zap.Error(err),
			zap.String("original_url", url),
		)
		return generateShortURL(*cfg.BaseURL, shortURLID), err
	} else if err != nil {
		s.Logger.Warn("unexpected error", zap.Error(err))
		return "", err
	}

	// s.EP.IncrementUUID()
	// event := eventlog.Event{
	// 	UUID:        strconv.Itoa(s.EP.CurrentUUID),
	// 	ShortURL:    ShortURLID,
	// 	OriginalURL: url,
	// }
	// s.EP.WriteEvent(&event)

	return fmt.Sprintf("%s/%s", *cfg.BaseURL, shortURLID), nil
}

func (s *ShortenService) ShortenURLBatch(
	cfg config.Config,
	batchData []dto.BatchShortenRequest,
) ([]dto.BatchShortenResponse, error) {
	s.Logger.Info("processing batch of urls")

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
		s.Logger.Warn("error while saving batch of urls", zap.Error(err))
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

func generateShortURL(baseURL, shortURLID string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortURLID)
}
