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
	SetShortURL(ctx context.Context, token string, url string) (string, error)
	BatchInsertShortURLS(ctx context.Context, entries []models.BatchInsertURLEntry) error
	GetNumberOfEntries(ctx context.Context) int
	PingConnect(ctx context.Context) error
}
