package repository

import (
	"go_link_shortener/internal/models"
)

type URLKeeper interface {
	FindByID(id string) string
	FindByLink(link string) string
	SetShortURL(encodedString string, url string)
	BatchInsertShortURLS(entries []models.BatchInsertURLEntry) error
	GetNumberOfEntries() int
	PingConnect() error
}
