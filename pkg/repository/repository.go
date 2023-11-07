package repository

import (
	"database/sql"
	"sync"
)

type URLRepository struct {
	URLEntries map[string]string
	lock       *sync.Mutex
	DB         *sql.DB
}

func NewRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		URLEntries: make(map[string]string),
		lock:       &sync.Mutex{},
		DB:         db,
	}
}

func (ur *URLRepository) FindByID(id string) (string, bool) {
	ur.lock.Lock()
	defer ur.lock.Unlock()
	entry, ok := ur.URLEntries[id]
	return entry, ok
}

func (ur *URLRepository) FindByLink(link string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	for index, value := range ur.URLEntries {
		if link == value {
			return index
		}
	}

	return ""
}

func (ur *URLRepository) SetShortURL(shortURL string, origURL string) string {
	ur.lock.Lock()
	defer ur.lock.Unlock()

	ur.URLEntries[shortURL] = origURL
	return shortURL
}

func (ur *URLRepository) GetShortURL(shortURL string) (string, bool) {
	entry, ok := ur.FindByID(shortURL)
	return entry, ok
}

func (ur *URLRepository) GetNumberOfEntries() int {
	return len(ur.URLEntries)
}

func (ur *URLRepository) PingDBConnect() error {
	err := ur.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
