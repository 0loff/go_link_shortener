package inmemoryrepository

import (
	"context"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/repository"
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

func (ur *InMemoryRepository) FindByID(ctx context.Context, id string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()
	entry, ok := ur.URLEntries[id]

	if !ok {
		return ""
	}
	return entry
}

func (ur *InMemoryRepository) FindByLink(ctx context.Context, link string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for index, value := range ur.URLEntries {
		if link == value {
			return index
		}
	}

	return ""
}

func (ur *InMemoryRepository) SetShortURL(ctx context.Context, shortURL string, origURL string) (string, error) {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for _, value := range ur.URLEntries {
		if origURL == value {
			return shortURL, repository.ErrConflict
		}
	}

	ur.URLEntries[shortURL] = origURL
	return shortURL, nil
}

func (ur *InMemoryRepository) BatchInsertShortURLS(ctx context.Context, urls []models.BatchInsertURLEntry) error {
	for _, u := range urls {
		ur.SetShortURL(ctx, u.ShortURL, u.OriginalURL)
	}
	return nil
}

func (ur *InMemoryRepository) GetNumberOfEntries(ctx context.Context) int {
	return len(ur.URLEntries)
}

func (ur *InMemoryRepository) PingConnect(ctx context.Context) error {
	return nil
}
