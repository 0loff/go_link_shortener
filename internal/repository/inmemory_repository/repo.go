package inmemoryrepository

import (
	"context"
	"sync"

	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/internal/repository"
)

// Структура записи сокращенного URL для хранения в slice
type InmemoryEntry struct {
	UserID      string
	ShortURL    string
	OriginalURL string
	IsDeleted   bool
}

// Структура репозитория для хранения сокращенных URL в памяти
type InMemoryRepository struct {
	URLEntries []InmemoryEntry
	lock       *sync.Mutex
}

// Инициализация репозитория для хранения записей в памяти сокращенных URL
func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		URLEntries: []InmemoryEntry{},
		lock:       &sync.Mutex{},
	}
}

// Поиск записи сокращенного URL по сгенерированному токену в slice
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

// Поиск записи сокращенного URL по оригинальному URL адресу
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

// Поиск всех записей сокращенных URL по пользователю
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

// Создание записи сокращенного URL в slice
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

// Создание множественной записи сокращенных URLs переданных пользователем в одном запросе
func (ur *InMemoryRepository) BatchInsertShortURLS(ctx context.Context, uid string, urls []models.URLEntry) error {
	for _, u := range urls {
		ur.SetShortURL(ctx, uid, u.ShortURL, u.OriginalURL)
	}
	return nil
}

// Установка флага удаления записи сокращенного URL
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

// Получение количества записей всех сокращенных URLs в slice
func (ur *InMemoryRepository) GetNumberOfEntries(ctx context.Context) int {
	return len(ur.URLEntries)
}

// GetMetrics method to get statistics about saved short urls and active users
func (ur *InMemoryRepository) GetMetrics() models.Metrics {
	users := make(map[string]struct{})

	for _, entry := range ur.URLEntries {
		if _, ok := users[entry.UserID]; !ok {
			users[entry.UserID] = struct{}{}
		}
	}

	return models.Metrics{
		Urls:  len(ur.URLEntries),
		Users: len(users),
	}
}

// Мокированный метод проверки соединения с файлом для имплементации интерфейса URLKeeper
func (ur *InMemoryRepository) PingConnect(ctx context.Context) error {
	return nil
}
