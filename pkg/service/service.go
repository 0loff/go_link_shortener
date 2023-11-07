package service

import (
	"fmt"
	filehandler "go_link_shortener/internal/fileHandler"
	"go_link_shortener/pkg/base62"
	"go_link_shortener/pkg/repository"
	"log"
)

type Service struct {
	Repo         repository.URLKeeper
	ShortURLHost string
	StorageFile  string
}

func NewService(Repo repository.URLKeeper, shortURLHost string, storageFile string) *Service {
	return &Service{
		Repo:         Repo,
		ShortURLHost: shortURLHost,
		StorageFile:  storageFile,
	}
}

func (s *Service) SetShortURL(url string) string {
	shortURL := s.Repo.SetShortURL(base62.NewBase62Encoder().EncodeString(), url)

	newEntry := filehandler.Entry{
		ID:          s.Repo.GetNumberOfEntries(),
		ShortURL:    shortURL,
		OriginalURL: url,
	}

	s.WriteToFile(newEntry)
	return shortURL
}

func (s *Service) GetShortURL(shortURL string) string {
	link, ok := s.Repo.FindByID(shortURL)
	if ok {
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

func (s *Service) WriteToFile(entry filehandler.Entry) {
	if s.StorageFile == "" {
		return
	}

	Producer, err := filehandler.NewProducer(s.StorageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()
	Producer.WriteEntry(&entry)
}

func (s *Service) RecoverFromFile() {
	if s.StorageFile == "" {
		return
	}

	Consumer, err := filehandler.NewConsumer(s.StorageFile)
	if err != nil {
		log.Fatal(err)
	}

	defer Consumer.Close()

	for {
		entry, err := Consumer.ReadEntry()
		if err != nil {
			break
		}

		if s.Repo.FindByLink(entry.OriginalURL) == "" {
			s.Repo.SetShortURL(entry.ShortURL, entry.OriginalURL)
			fmt.Println("Short url: ", entry.ShortURL, "-", entry.OriginalURL, " has been recovered from file.")
		}
	}
}

func (s *Service) SetTestShortURL() {
	shortURLforTest := s.Repo.FindByLink("https://practicum.yandex.ru/")
	if shortURLforTest == "" {
		s.Repo.SetShortURL("OL0ZGlVC3dq", "https://practicum.yandex.ru/")
	}
}
