package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/mocks"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestDeleteURLsBatch(t *testing.T) {
	t.Run("test delete batch", func(t *testing.T) {
		r := &inmemory.Memory{}
		err := r.Init(context.Background(), cfg, l)

		assert.NoError(t, err)

		s := initTestService(r)
		expected := model.BatchDeleteShortURLs{
			ShortURLs: []string{"abc123"},
			UserID:    uuid.NewString(),
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		s.DeleteURLsBatch(ctx, expected)

		select {
		case actual := <-s.CommCh:
			assert.Equal(t, expected, actual)
		case <-time.After(time.Second):
			t.Fatal("timeout: value was not written to channel")
		}
	})
}

func TestProcessDeletions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testData := model.BatchDeleteShortURLs{
		UserID:    "user42",
		ShortURLs: []string{"abc", "xyz"},
	}

	mockStorage := mocks.NewMockRepository(ctrl)

	mockStorage.EXPECT().
		DeleteBatch(gomock.Any(), testData).
		Return(int64(2), nil).Times(1)

	s := initTestService(mockStorage)

	done := make(chan struct{})

	go func() {
		s.ProcessDeletions()
		close(done)
	}()

	s.CommCh <- testData
	close(s.CommCh)

	select {
	case <-done:
	case <-time.After(time.Second * 30):
		t.Fatal("timeout waiting for ProcessDeletions to finish")
	}

}
