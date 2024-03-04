package repository

import (
	"context"
	"errors"

	"github.com/0loff/go_link_shortener/internal/models"
)

// Кастомные ошибки взаимодействия с БД при запросах на получение сохраненных URLs
var (
	ErrConflict    = errors.New("data conflict")
	ErrURLGone     = errors.New("URL have been deleted")
	ErrURLNotFound = errors.New("URL not found")
)

// Интерфейс взаимодействия сервисов с хранилищем сокращенных URLs
type URLKeeper interface {
	FindByID(ctx context.Context, id string) (string, error)
	FindByLink(ctx context.Context, link string) string
	FindByUser(ctx context.Context, uid string) []models.URLEntry
	SetShortURL(ctx context.Context, uid, token, url string) (string, error)
	BatchInsertShortURLS(ctx context.Context, uid string, entries []models.URLEntry) error
	SetDelShortURLS(URLSList []models.DelURLEntry) error
	GetNumberOfEntries(ctx context.Context) int
	GetMetrics() models.Metrics
	PingConnect(ctx context.Context) error
}
