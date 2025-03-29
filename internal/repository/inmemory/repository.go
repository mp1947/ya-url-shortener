package inmemory

import "github.com/mp1947/ya-url-shortener/config"

type Memory struct {
	data        map[string]string
	StorageType string
}

func (s *Memory) Init(cfg config.Config) error {
	s.data = make(map[string]string)
	s.StorageType = "inmemory"
	return nil
}

func (s *Memory) Save(shortURL, originalURL string) bool {
	isSaved := false
	if s.data[shortURL] == "" {
		s.data[shortURL] = originalURL
		isSaved = true
	}
	return isSaved
}

func (s *Memory) Get(shortURL string) string {
	return s.data[shortURL]
}

func (s *Memory) GetType() string {
	return s.StorageType
}

func (s *Memory) Ping() error {
	return nil
}
