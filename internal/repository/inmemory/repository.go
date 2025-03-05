package inmemory

type Memory struct {
	data map[string]string
}

func Init() *Memory {
	return &Memory{data: make(map[string]string)}
}

func (s *Memory) Save(shortURL, originalURL string) {
	s.data[shortURL] = originalURL
}

func (s *Memory) Get(shortURL string) string {
	return s.data[shortURL]
}
