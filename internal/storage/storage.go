package storage

type URLStorage struct {
	urls map[string]string
}

func New() *URLStorage {
	return &URLStorage{
		urls: make(map[string]string),
	}
}

func (s *URLStorage) SaveUrl(url string, alias string) {
	s.urls[alias] = url
}

func (s *URLStorage) GetUrl(alias string) (url string, exists bool) {
	url, exists = s.urls[alias]
	return url, exists
}
