package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/mocks"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserURLS(t *testing.T) {
	t.Run("test get user urls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockRepository(ctrl)

		testUserID := uuid.NewString()

		mockStorage.EXPECT().
			GetURLsByUserID(gomock.Any(), testUserID).
			Return([]model.UserURL{
				{
					ShortURLID:  "aaabbb",
					OriginalURL: "https://google.com",
				},
				{
					ShortURLID:  "eeebasbdh",
					OriginalURL: "https://yandex.com",
				},
			}, nil).Times(1)

		s := initTestService(mockStorage)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

		defer cancel()

		urls, err := s.GetUserURLs(ctx, testUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, urls)
	})
}
