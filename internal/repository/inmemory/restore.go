package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"go.uber.org/zap"
)

// RestoreFromFile restores the in-memory storage state from a file specified in the configuration.
// It reads each line from the file, unmarshals it into an eventlog.Event, and saves it to the storage.
// The method returns the number of records restored and any error encountered during the process.
// If an error occurs while saving a record, it logs a warning but continues processing the rest of the file.
func (s *Memory) RestoreFromFile(l *zap.Logger) (int, error) {
	file, err := os.OpenFile(*s.cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	currentUUID := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.isInRestoreMode = true

	for scanner.Scan() {
		var event eventlog.Event
		line := scanner.Text()

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return 0, err
		}
		if err := s.Save(ctx, event.ShortURL, event.OriginalURL, event.UserID); err != nil {
			l.Warn("error saving record to file during restore phase", zap.Error(err))
		}

		currentUUID += 1
	}
	s.EP.CurrentUUID = currentUUID
	s.isInRestoreMode = false

	return currentUUID, nil
}
