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

func (s *Service) CreateShortURL(ctx context.Context, uid, url string) (string, error) {
	token := base62.NewBase62Encoder().EncodeString()

	shortURL, err := s.Repo.SetShortURL(ctx, uid, token, url)
	if err != nil && errors.Is(err, repository.ErrConflict) {
		shortURL = s.Repo.FindByLink(ctx, url)
		logger.Log.Error("Error when inserting short URL into database", zap.Error(err))
	}

	return shortURL, err
}

func (s *Service) SetBatchShortURLs(ctx context.Context, uid string, entries []models.BatchURLRequestEntry) []models.BatchURLResponseEntry {
	batchEntries := []models.URLEntry{}
	respEntries := []models.BatchURLResponseEntry{}

	for _, u := range entries {
		shortURL := s.Repo.FindByLink(ctx, u.OriginalURL)

		if shortURL == "" {
			shortURL = base62.NewBase62Encoder().EncodeString()

			batchEntries = append(batchEntries, models.URLEntry{
				ShortURL:    shortURL,
				OriginalURL: u.OriginalURL,
			})
		}

		respEntries = append(respEntries, models.BatchURLResponseEntry{
			CorrelationID: u.CorrelationID,
			ShortURL:      s.ShortURLHost + "/" + shortURL,
		})
	}

	if len(batchEntries) != 0 {
		err := s.Repo.BatchInsertShortURLS(ctx, uid, batchEntries)
		if err != nil {
			logger.Log.Error("BatchInsert failed", zap.Error(err))
		}
	}

	return respEntries
}

func (s *Service) GetShortURL(ctx context.Context, shortURL string) string {
	return s.Repo.FindByID(ctx, shortURL)
}

func (s *Service) GetShortURLs(ctx context.Context, uid string) []models.URLEntry {
	var UserURLs []models.URLEntry

	userURLs := s.Repo.FindByUser(ctx, uid)
	for _, u := range userURLs {
		UserURLs = append(UserURLs, models.URLEntry{
			ShortURL:    s.ShortURLHost + "/" + u.ShortURL,
			OriginalURL: u.OriginalURL,
		})
	}

	return UserURLs
}
