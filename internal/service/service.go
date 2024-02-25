package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/internal/repository"
	"github.com/0loff/go_link_shortener/pkg/base62"
)

// Структура для инициализации сервиса
type Service struct {
	Repo         repository.URLKeeper
	ShortURLHost string
	StorageFile  string
	DelCh        chan models.DelURLEntry
}

// Конструктор инициализации сервиса при запуске приложения
func NewService(Repo repository.URLKeeper, shortURLHost string) *Service {
	service := &Service{
		Repo:         Repo,
		ShortURLHost: shortURLHost,
		DelCh:        make(chan models.DelURLEntry, 1024),
	}

	go service.DeleteManager(service.DelCh)

	return service
}

// Метод создания сокращенного URL с последующим вызовом репозитория для сохранения
func (s *Service) CreateShortURL(ctx context.Context, uid, url string) (string, error) {
	seed := time.Now().UnixNano()
	token := base62.NewBase62Encoder().EncodeString(uint64(seed))

	shortURL, err := s.Repo.SetShortURL(ctx, uid, token, url)
	if err != nil && errors.Is(err, repository.ErrConflict) {
		shortURL = s.Repo.FindByLink(ctx, url)
		logger.Log.Error("Error when inserting short URL into database", zap.Error(err))
	}

	return shortURL, err
}

// Метод создания множества записей сокращенных URLs одним пользовательским запросом
func (s *Service) SetBatchShortURLs(ctx context.Context, uid string, entries []models.BatchURLRequestEntry) []models.BatchURLResponseEntry {
	batchEntries := []models.URLEntry{}
	respEntries := []models.BatchURLResponseEntry{}

	seed := time.Now().UnixNano()

	for _, u := range entries {
		shortURL := s.Repo.FindByLink(ctx, u.OriginalURL)

		if shortURL == "" {
			shortURL = base62.NewBase62Encoder().EncodeString(uint64(seed))

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

// Метод получения сокращенного URL
func (s *Service) GetShortURL(ctx context.Context, shortURL string) (string, error) {
	test, err := s.Repo.FindByID(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return test, nil
}

// Метод полуения нескольких записей сокращенных URLs по текущему пользователю
func (s *Service) GetShortURLs(ctx context.Context, uid string) []models.URLEntry {
	var UserURLs []models.URLEntry

	userURLs := s.Repo.FindByUser(ctx, uid)
	for _, u := range userURLs {
		if u.IsDeleted {
			continue
		}

		UserURLs = append(UserURLs, models.URLEntry{
			ShortURL:    s.ShortURLHost + "/" + u.ShortURL,
			OriginalURL: u.OriginalURL,
		})
	}

	return UserURLs
}

// Метод установки флага удаления сокращенных URLs, переданных списком в одном запросе
func (s *Service) DelShortURLs(uid string, URLSList []string) {
	var URLEnties []models.DelURLEntry

	for _, URL := range URLSList {
		URLEnties = append(URLEnties, models.DelURLEntry{
			UserID:   uid,
			ShortURL: URL,
		})
	}

	ShortURLSCh := s.ChGenerator(URLEnties)
	s.MergeChs(ShortURLSCh)
}

// Горутина, выполняющая запрос на удаление списка сокращенных урлов по таймингу (один раз в 10 секунд)
func (s *Service) DeleteManager(URLCh chan models.DelURLEntry) {
	ticker := time.NewTicker(10 * time.Second)

	var URLSForDel []models.DelURLEntry

	for {
		select {
		case ShortURL := <-URLCh:
			URLSForDel = append(URLSForDel, ShortURL)

		case <-ticker.C:
			if len(URLSForDel) == 0 {
				continue
			}
			s.Repo.SetDelShortURLS(URLSForDel)
			URLSForDel = nil
		}
	}
}

// Создание канала для постановки в очередь на удаление сокращенных URLs
func (s *Service) ChGenerator(ShortURLSlist []models.DelURLEntry) chan models.DelURLEntry {
	inputCh := make(chan models.DelURLEntry)

	go func() {
		defer close(inputCh)

		for _, URLEntry := range ShortURLSlist {
			inputCh <- URLEntry
		}
	}()

	return inputCh
}

// Реализация паттерна fanOut для мержа нескольких каналов на удаление в единую очередь для DeleteManager'a
func (s *Service) MergeChs(resultChan ...chan models.DelURLEntry) {
	var wg sync.WaitGroup

	for _, ch := range resultChan {
		chClosure := ch
		wg.Add(1)

		go func() {
			for data := range chClosure {
				s.DelCh <- data
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
	}()
}
