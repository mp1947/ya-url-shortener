package inmemory

type Storage struct {
	data map[string]string
}

func InitStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

func (s *Storage) Save(shortURL, originalURL string) {
	s.data[shortURL] = originalURL
}

func (s *Storage) Get(shortURL string) string {
	return s.data[shortURL]
}
