package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/mp1947/ya-url-shortener/internal/mocks"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetOriginalURL(t *testing.T) {
	t.Run("test get url", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStorage := mocks.NewMockRepository(ctrl)

		testDataExpected := "https://whatever.com"
		testDataToSend := usecase.GenerateIDFromURL(testDataExpected)

		mockStorage.EXPECT().
			Get(gomock.Any(), testDataToSend).
			Return(model.URL{
				ShortURLID:  testDataExpected,
				OriginalURL: "https://whatever.com",
				IsDeleted:   false,
			}, nil).Times(1)

		s := initTestService(mockStorage)

		ctx, shutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer shutdown()

		url, err := s.GetOriginalURL(ctx, testDataToSend)

		assert.NoError(t, err)
		assert.Equal(t, testDataExpected, url.OriginalURL)
	})
}
