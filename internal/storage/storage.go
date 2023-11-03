package storage

import (
	"fmt"
	"log"
	"sync"
)

var Store LinkStorageRepository

type LinkStorageRepository struct {
	mu            sync.Mutex
	linkEntries   map[string]string
	shortLinkHost string
	storageFile   string
}

func (ls *LinkStorageRepository) LinkStorageCreate() {
	ls.linkEntries = make(map[string]string)
}

func (ls *LinkStorageRepository) LinkStorageRecover() {
	if ls.GetStorageFile() == "" {
		return
	}

	Consumer, err := NewConsumer(ls.GetStorageFile())
	if err != nil {
		log.Fatal(err)
	}

	defer Consumer.Close()

	for {
		entry, err := Consumer.ReadEntry()
		if err != nil {
			break
		}
		ls.linkEntries[entry.ShortURL] = entry.OriginalURL
		fmt.Println("Short url: ", entry.ShortURL, "-", entry.OriginalURL, " has been recovered from file.")
	}

	testShortURL := Store.FindByLink("https://practicum.yandex.ru/")
	if testShortURL == "" {
		Store.SetShortLink("OL0ZGlVC3dq", "https://practicum.yandex.ru/")
	}
}

func (ls *LinkStorageRepository) FindByID(id string) (string, bool) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	entry, ok := ls.linkEntries[id]
	return entry, ok
}

func (ls *LinkStorageRepository) FindByLink(link string) string {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for index, value := range ls.linkEntries {
		if link == value {
			return index
		}
	}

	return ""
}

func (ls *LinkStorageRepository) SetShortLink(shortURL string, originURL string) string {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.linkEntries[shortURL] = originURL

	entry := Entry{
		ID:          len(Store.linkEntries),
		ShortURL:    shortURL,
		OriginalURL: originURL,
	}

	ls.SaveEntryToFile(entry)

	return shortURL
}

func (ls *LinkStorageRepository) GetShortLink(shortURL string) (string, bool) {
	entry, ok := ls.FindByID(shortURL)
	return entry, ok
}

func (ls *LinkStorageRepository) SetShortLinkHost(shortLinkHost string) {
	ls.shortLinkHost = shortLinkHost
}

func (ls *LinkStorageRepository) GetShortLinkHost() string {
	return ls.shortLinkHost
}

func (ls *LinkStorageRepository) SetStorageFile(storageFile string) {
	ls.storageFile = storageFile
}

func (ls *LinkStorageRepository) GetStorageFile() string {
	return ls.storageFile
}

func (ls *LinkStorageRepository) SaveEntryToFile(entry Entry) {
	if ls.GetStorageFile() == "" {
		return
	}

	Producer, err := NewProducer(ls.GetStorageFile())
	if err != nil {
		log.Fatal(err)
	}

	defer Producer.Close()

	Producer.WriteEntry(&entry)
}

func LinkStorageInit() {
	Store = LinkStorageRepository{}
	Store.LinkStorageCreate()
}
