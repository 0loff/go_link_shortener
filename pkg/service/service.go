package service

import (
	"context"
	"errors"
	"go_link_shortener/internal/logger"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/base62"
	"go_link_shortener/pkg/repository"

	"go.uber.org/zap"
)

type Service struct {
	Repo         repository.URLKeeper
	ShortURLHost string
	StorageFile  string
}

func NewService(Repo repository.URLKeeper, shortURLHost string) *Service {
	return &Service{
		Repo:         Repo,
		ShortURLHost: shortURLHost,
	}
}

func (s *Service) CreateShortURL(ctx context.Context, url string) (string, error) {
	token := base62.NewBase62Encoder().EncodeString()

	shortURL, err := s.Repo.SetShortURL(ctx, token, url)
	if err != nil && errors.Is(err, repository.ErrConflict) {
		shortURL = s.Repo.FindByLink(ctx, url)
		logger.Log.Error("Error when inserting short URL into database", zap.Error(err))
	}

	return shortURL, err
}

func (s *Service) SetBatchShortURLs(ctx context.Context, entries []models.BatchURLRequestEntry) []models.BatchURLResponseEntry {
	batchEntries := []models.BatchInsertURLEntry{}
	respEntries := []models.BatchURLResponseEntry{}

	for _, u := range entries {
		shortURL := s.Repo.FindByLink(ctx, u.OriginalURL)

		if shortURL == "" {
			shortURL = base62.NewBase62Encoder().EncodeString()
			newInsertEntry := models.BatchInsertURLEntry{
				ShortURL:    shortURL,
				OriginalURL: u.OriginalURL,
			}
			batchEntries = append(batchEntries, newInsertEntry)
		}

		newResponseEntry := models.BatchURLResponseEntry{
			CorrelationID: u.CorrelationID,
			ShortURL:      s.ShortURLHost + "/" + shortURL,
		}
		respEntries = append(respEntries, newResponseEntry)
	}
	if len(batchEntries) != 0 {
		err := s.Repo.BatchInsertShortURLS(ctx, batchEntries)
		if err != nil {
			logger.Log.Error("BatchInsert failed", zap.Error(err))
		}
	}
	return respEntries
}

func (s *Service) GetShortURL(ctx context.Context, shortURL string) string {
	link := s.Repo.FindByID(ctx, shortURL)

	if link != "" {
		return link
	}

	return ""
}
