package repository

import (
	"context"
	"errors"
	"go_link_shortener/internal/models"
)

var ErrConflict = errors.New("data conflict")

type URLKeeper interface {
	FindByID(ctx context.Context, id string) string
	FindByLink(ctx context.Context, link string) string
	FindByUser(ctx context.Context, uid string) []models.URLEntry
	SetShortURL(ctx context.Context, uid, token, url string) (string, error)
	BatchInsertShortURLS(ctx context.Context, uid string, entries []models.URLEntry) error
	GetNumberOfEntries(ctx context.Context) int
	PingConnect(ctx context.Context) error
}
