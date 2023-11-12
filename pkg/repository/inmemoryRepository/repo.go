package inmemoryrepository

import (
	"go_link_shortener/internal/models"
	"sync"
)

type InMemoryRepository struct {
	URLEntries map[string]string
	lock       *sync.Mutex
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		URLEntries: make(map[string]string),
		lock:       &sync.Mutex{},
	}
}

func (ur *InMemoryRepository) FindByID(id string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()
	entry, ok := ur.URLEntries[id]

	if !ok {
		return ""
	}
	return entry
}

func (ur *InMemoryRepository) FindByLink(link string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for index, value := range ur.URLEntries {
		if link == value {
			return index
		}
	}

	return ""
}

func (ur *InMemoryRepository) SetShortURL(shortURL string, origURL string) {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	ur.URLEntries[shortURL] = origURL
}

func (ur *InMemoryRepository) BatchInsertShortURLS(urls []models.BatchInsertURLEntry) error {
	for _, u := range urls {
		ur.SetShortURL(u.ShortURL, u.OriginalURL)
	}
	return nil
}

func (ur *InMemoryRepository) GetNumberOfEntries() int {
	return len(ur.URLEntries)
}

func (ur *InMemoryRepository) PingConnect() error {
	return nil
}
