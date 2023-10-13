package main

import (
	"log"
)

var storage LinkStorageRepository

type LinkStorageRepository struct {
	linkEntries map[string]string
}

func (ls LinkStorageRepository) LinkStorageCreate() LinkStorageRepository {
	ls.linkEntries = make(map[string]string)
	return ls
}

func (ls LinkStorageRepository) FindByID(id string) (string, bool) {
	entry, ok := ls.linkEntries[id]
	return entry, ok
}

func (ls LinkStorageRepository) FindByLink(link string) string {
	for index, value := range ls.linkEntries {
		if link == value {
			return index
		}
	}

	return ""
}

func (ls LinkStorageRepository) SetShortLink(shortURL string, originURL string) LinkStorageRepository {
	ls.linkEntries[shortURL] = originURL
	return ls
}

func (ls LinkStorageRepository) GetShortLink(shortURL string) (string, bool) {
	entry, ok := ls.FindByID(shortURL)
	return entry, ok
}

func (ls LinkStorageRepository) LogAllEntries() {
	log.Println(ls)
}

func LinkStorageInit() {
	storage = new(LinkStorageRepository).LinkStorageCreate()
	storage.SetShortLink("OL0ZGlVC3dq", "https://practicum.yandex.ru/")
}
