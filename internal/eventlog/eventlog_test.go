package eventlog_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/stretchr/testify/assert"
)

func TestEventProcessor(t *testing.T) {
	cfg := &config.Config{}
	cfg.InitConfig()

	ep, err := eventlog.NewEventProcessor(*cfg)

	t.Run("test event processor creation", func(t *testing.T) {
		assert.NoError(t, err)
		assert.NotNil(t, ep)
	})

	t.Run("test event write", func(t *testing.T) {
		err = ep.WriteEvent(&eventlog.Event{
			UUID:        uuid.NewString(),
			ShortURL:    "testtest",
			OriginalURL: "testtest",
			UserID:      "someUserID",
			IsDeleted:   false,
		})
		assert.NoError(t, err)
	})

	t.Run("test id increment", func(t *testing.T) {
		assert.NoError(t, err)
		before := ep.CurrentUUID
		ep.IncrementUUID()
		after := ep.CurrentUUID
		assert.Greater(t, after, before)
	})
}
