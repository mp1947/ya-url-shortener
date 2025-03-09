package inmemory

type Memory struct {
	data map[string]string
}

func (s *Memory) Init() {
	s.data = make(map[string]string)
}

func (s *Memory) Save(shortURL, originalURL string) {
	s.data[shortURL] = originalURL
}

func (s *Memory) Get(shortURL string) string {
	return s.data[shortURL]
}
