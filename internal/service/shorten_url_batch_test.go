package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	"github.com/mp1947/ya-url-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShortenURLBatch(t *testing.T) {
	t.Run("shorten batch of the urls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		correlationID := "aaabbbccc"

		urls := []dto.BatchShortenRequest{
			{
				CorrelationID: correlationID,
				OriginalURL:   "https://google.com",
			},
			{
				CorrelationID: correlationID,
				OriginalURL:   "https://yandex.com",
			},
		}

		userID := uuid.NewString()

		mockRepository := mocks.NewMockRepository(ctrl)

		mockRepository.EXPECT().
			SaveBatch(gomock.Any(), gomock.Any(), userID).
			Return(true, nil).Times(1)

		s := initTestService(mockRepository)

		ctx, shutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer shutdown()

		out, err := s.ShortenURLBatch(ctx, urls, userID)
		assert.NoError(t, err)
		assert.Len(t, out, len(urls))
		assert.IsType(t, []dto.BatchShortenResponse{}, out)
	})
}
