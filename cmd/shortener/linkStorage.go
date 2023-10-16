package main

import (
	"log"
	"sync"
)

var storage LinkStorageRepository

type LinkStorageRepository struct {
	mu          sync.Mutex
	linkEntries map[string]string
}

func (ls *LinkStorageRepository) LinkStorageCreate() {
	ls.linkEntries = make(map[string]string)
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

func (ls *LinkStorageRepository) SetShortLink(shortURL string, originURL string) {
	ls.mu.Lock()
	ls.linkEntries[shortURL] = originURL
	ls.mu.Unlock()
}

func (ls *LinkStorageRepository) GetShortLink(shortURL string) (string, bool) {
	entry, ok := ls.FindByID(shortURL)
	return entry, ok
}

func (ls *LinkStorageRepository) LogAllEntries() {
	log.Println(ls)
}

func LinkStorageInit() {
	storage = LinkStorageRepository{}
	storage.LinkStorageCreate()
	storage.SetShortLink("OL0ZGlVC3dq", "https://practicum.yandex.ru/")
}
