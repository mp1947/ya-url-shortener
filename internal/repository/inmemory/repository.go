package inmemory

import (
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/entity"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

type Memory struct {
	data        map[string]string
	StorageType string
}

func (s *Memory) Init(cfg config.Config) error {
	s.data = make(map[string]string)
	s.StorageType = "inmemory"
	return nil
}

func (s *Memory) Save(shortURL, originalURL string) error {
	if s.data[shortURL] == "" {
		s.data[shortURL] = originalURL
		return nil
	}
	return shrterr.ErrOriginalURLAlreadyExists
}

func (s *Memory) SaveBatch(urls []entity.URL) (bool, error) {
	for _, v := range urls {
		s.data[v.ShortURLID] = v.OriginalURL
	}
	return true, nil
}

func (s *Memory) Get(shortURL string) (string, error) {
	return s.data[shortURL], nil
}

func (s *Memory) GetType() string {
	return s.StorageType
}

func (s *Memory) Ping() error {
	return nil
}
