package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/mocks"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShortenURL(t *testing.T) {
	t.Run("test shorten url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepository := mocks.NewMockRepository(ctrl)

		urlToTest := "https://google.com"
		shortURL := usecase.GenerateIDFromURL(urlToTest)
		userID := uuid.NewString()

		mockRepository.EXPECT().
			Save(gomock.Any(), shortURL, urlToTest, userID).
			Return(nil).Times(1)

		mockRepository.EXPECT().
			Save(gomock.Any(), shortURL, urlToTest, userID).
			Return(errors.New("some error")).Times(1)

		s := initTestService(mockRepository)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		out, err := s.ShortenURL(ctx, urlToTest, userID)

		assert.NoError(t, err)
		assert.Equal(t, baseURL+"/"+shortURL, out)

		newOut, err := s.ShortenURL(ctx, urlToTest, userID)

		assert.Error(t, err)
		assert.Empty(t, newOut)
	})
}
