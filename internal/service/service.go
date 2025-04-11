package service

import (
	"context"
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
	ShortenURL(
		ctx context.Context,
		cfg config.Config,
		url string,
		userID string,
	) (string, error)
	GetOriginalURL(ctx context.Context, shortURLID string) (string, error)
	ShortenURLBatch(
		ctx context.Context,
		cfg config.Config,
		batchData []dto.BatchShortenRequest,
		userID string,
	) ([]dto.BatchShortenResponse, error)
	GetUserURLs(
		ctx context.Context,
		cfg config.Config,
		userID string,
	) ([]dto.ShortenURLsByUserID, error)
}

type ShortenService struct {
	Storage repository.Repository
	EP      eventlog.EventProcessor
	Logger  *zap.Logger
}

func (s *ShortenService) ShortenURL(
	ctx context.Context,
	cfg config.Config,
	url string,
	userID string,
) (string, error) {
	s.Logger.Info("shortening incoming url", zap.String("original_url", url))

	shortURLID := usecase.GenerateIDFromURL(url)

	s.Logger.Info(
		"short_url id generated for url",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", url),
		zap.String("user_id", userID),
	)

	err := s.Storage.Save(ctx, shortURLID, url, userID)
	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		s.Logger.Info(
			"original_url already exists, returning error with short url",
			zap.Error(err),
			zap.String("original_url", url),
		)
		return generateShortURL(*cfg.BaseURL, shortURLID), err
	} else if err != nil {
		s.Logger.Warn("unexpected error", zap.Error(err))
		return "", err
	}

	return generateShortURL(*cfg.BaseURL, shortURLID), nil
}

func (s *ShortenService) ShortenURLBatch(
	ctx context.Context,
	cfg config.Config,
	batchData []dto.BatchShortenRequest,
	userID string,
) ([]dto.BatchShortenResponse, error) {
	s.Logger.Info(
		"processing batch of urls",
		zap.Any("batch_data", batchData),
	)

	urls := make([]entity.URLWithCorrelation, len(batchData))
	result := make([]dto.BatchShortenResponse, len(batchData))

	for i, v := range batchData {
		shortURLID := usecase.GenerateIDFromURL(v.OriginalURL)
		urls[i] = entity.URLWithCorrelation{
			ShortURLID:    shortURLID,
			OriginalURL:   v.OriginalURL,
			CorrelationID: v.CorrelationID,
		}
		result[i] = dto.BatchShortenResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", *cfg.BaseURL, shortURLID),
		}
	}

	_, err := s.Storage.SaveBatch(ctx, urls, userID)
	if err != nil {
		s.Logger.Warn("error while saving batch of urls", zap.Error(err))
		return nil, err
	}

	s.Logger.Info("batch of urls were successfully processed")

	return result, nil

}

func (s *ShortenService) GetOriginalURL(
	ctx context.Context,
	shortURLID string,
) (string, error) {

	s.Logger.Info("processing short url with id", zap.String("short_url_id", shortURLID))
	data, err := s.Storage.Get(ctx, shortURLID)
	if err != nil {
		s.Logger.Warn("error getting original_url by short_url_id", zap.String("short_url", shortURLID))
		return "", err
	}
	s.Logger.Info(
		"retrieved original_url by short_url_id",
		zap.String("short_url_id", shortURLID),
		zap.String("original_url", data),
	)
	return data, nil
}

func (s *ShortenService) GetUserURLs(
	ctx context.Context,
	cfg config.Config,
	userID string,
) ([]dto.ShortenURLsByUserID, error) {
	s.Logger.Info(
		"processing shorten urls from storage for user",
		zap.String("user_id", userID),
	)

	userURLs, err := s.Storage.GetURLsByUserID(ctx, userID)

	if err != nil {
		s.Logger.Warn("error getting urls by user id", zap.Error(err))
		return nil, err
	}

	userURLsResponse := make([]dto.ShortenURLsByUserID, len(userURLs))

	for i, v := range userURLs {
		userURLsResponse[i] = dto.ShortenURLsByUserID{
			ShortURL:    generateShortURL(*cfg.BaseURL, v.ShortURLID),
			OriginalURL: v.OriginalURL,
		}
	}
	return userURLsResponse, nil
}

func generateShortURL(baseURL, shortURLID string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortURLID)
}
