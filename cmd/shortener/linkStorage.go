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

func (ls LinkStorageRepository) FindById(id string) (string, bool) {
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

func (ls LinkStorageRepository) SetShortLink(shortUrl string, originUrl string) LinkStorageRepository {
	ls.linkEntries[shortUrl] = originUrl
	return ls
}

func (ls LinkStorageRepository) GetShortLink(shortUrl string) (string, bool) {
	entry, ok := ls.FindById(shortUrl)
	return entry, ok
}

func (ls LinkStorageRepository) LogAllEntries() {
	log.Println(ls)
}

func LinkStorageInit() {
	storage = new(LinkStorageRepository).LinkStorageCreate()
	storage.SetShortLink("OL0ZGlVC3dq", "https://practicum.yandex.ru/")
}
