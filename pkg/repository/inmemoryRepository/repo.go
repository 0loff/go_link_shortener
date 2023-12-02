package inmemoryrepository

import (
	"context"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/repository"
	"sync"
)

type InmemoryEntry struct {
	UserID      string
	ShortURL    string
	OriginalURL string
	IsDeleted   bool
}

type InMemoryRepository struct {
	URLEntries []InmemoryEntry
	lock       *sync.Mutex
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		URLEntries: []InmemoryEntry{},
		lock:       &sync.Mutex{},
	}
}

func (ur *InMemoryRepository) FindByID(ctx context.Context, id string) (string, error) {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for _, entry := range ur.URLEntries {
		if entry.ShortURL == id {
			if entry.IsDeleted {
				return "", repository.ErrURLGone
			}

			return entry.OriginalURL, nil
		}
	}

	return "", repository.ErrURLNotFound
}

func (ur *InMemoryRepository) FindByLink(ctx context.Context, link string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for _, entry := range ur.URLEntries {
		if entry.OriginalURL == link {
			return entry.ShortURL
		}
	}

	return ""
}

func (ur *InMemoryRepository) FindByUser(ctx context.Context, uid string) []models.URLEntry {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	var URLEntries []models.URLEntry

	for _, entry := range ur.URLEntries {
		if entry.UserID == uid && !entry.IsDeleted {
			URLEntries = append(URLEntries, models.URLEntry{
				ShortURL:    entry.ShortURL,
				OriginalURL: entry.OriginalURL,
			})
		}
	}

	return URLEntries
}

func (ur *InMemoryRepository) SetShortURL(ctx context.Context, uid, shortURL, origURL string) (string, error) {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for _, entry := range ur.URLEntries {
		if origURL == entry.OriginalURL {
			return shortURL, repository.ErrConflict
		}
	}

	ur.URLEntries = append(ur.URLEntries, InmemoryEntry{
		UserID:      uid,
		ShortURL:    shortURL,
		OriginalURL: origURL,
	})

	return shortURL, nil
}

func (ur *InMemoryRepository) BatchInsertShortURLS(ctx context.Context, uid string, urls []models.URLEntry) error {
	for _, u := range urls {
		ur.SetShortURL(ctx, uid, u.ShortURL, u.OriginalURL)
	}
	return nil
}

func (ur *InMemoryRepository) SetDelShortURLS(ShortURLsList []models.DelURLEntry) error {
	var UpdatedEntries []InmemoryEntry
	for _, entry := range ur.URLEntries {
		for _, URLForDel := range ShortURLsList {
			if entry.ShortURL == URLForDel.ShortURL && entry.UserID == URLForDel.UserID {
				entry.IsDeleted = true
			}
		}

		UpdatedEntries = append(UpdatedEntries, entry)
	}

	ur.URLEntries = UpdatedEntries
	return nil
}

func (ur *InMemoryRepository) GetNumberOfEntries(ctx context.Context) int {
	return len(ur.URLEntries)
}

func (ur *InMemoryRepository) PingConnect(ctx context.Context) error {
	return nil
}
