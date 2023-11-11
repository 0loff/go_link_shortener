package service

import (
	"go_link_shortener/pkg/base62"
	"go_link_shortener/pkg/repository"
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

func (s *Service) SetShortURL(url string) string {
	shortURL := base62.NewBase62Encoder().EncodeString()
	s.Repo.SetShortURL(shortURL, url)
	return shortURL
}

func (s *Service) GetShortURL(shortURL string) string {
	link := s.Repo.FindByID(shortURL)

	if link != "" {
		return link
	}

	return ""
}

func (s *Service) ShortURLResolver(url string) string {
	shortURL := s.Repo.FindByLink(url)

	if shortURL != "" {
		return shortURL
	}

	return s.SetShortURL(url)
}
