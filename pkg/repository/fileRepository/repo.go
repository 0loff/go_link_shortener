package filerepository

import (
	filehandler "go_link_shortener/internal/fileHandler"
	"go_link_shortener/internal/models"
	"log"
)

type FileRepository struct {
	StorageFile string
}

func NewRepository(fileName string) *FileRepository {
	return &FileRepository{
		StorageFile: fileName,
	}
}

func (fr *FileRepository) FindByID(id string) string {

	Consumer, err := filehandler.NewConsumer(fr.StorageFile)
	if err != nil {
		log.Fatal(err)
	}

	defer Consumer.Close()

	for {
		entry, err := Consumer.ReadEntry()
		// TODO err == io.EOF
		if err != nil {
			return ""
		}

		if entry.ShortURL == id {
			return entry.OriginalURL
		}
	}

}

func (fr *FileRepository) FindByLink(link string) string {
	Consumer, err := filehandler.NewConsumer(fr.StorageFile)
	if err != nil {
		log.Fatal(err)
	}

	defer Consumer.Close()

	for {
		entry, err := Consumer.ReadEntry()
		if err != nil {
			return ""
		}

		if entry.OriginalURL == link {
			return entry.ShortURL
		}
	}
}

func (fr *FileRepository) SetShortURL(shortURL, originURL string) {
	newEntry := filehandler.Entry{
		ID:          fr.GetNumberOfEntries(),
		ShortURL:    shortURL,
		OriginalURL: originURL,
	}

	fr.WriteToFile(newEntry)
}

func (fr *FileRepository) BatchInsertShortURLS(urls []models.BatchInsertURLEntry) error {
	for _, u := range urls {
		fr.WriteToFile(filehandler.Entry{
			ID:          fr.GetNumberOfEntries(),
			ShortURL:    u.ShortURL,
			OriginalURL: u.OriginalURL,
		})
	}

	return nil
}

func (fr *FileRepository) WriteToFile(entry filehandler.Entry) {
	Producer, err := filehandler.NewProducer(fr.StorageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()

	Producer.WriteEntry(&entry)
}

func (fr *FileRepository) GetNumberOfEntries() int {
	NumEntries := 0
	Consumer, err := filehandler.NewConsumer(fr.StorageFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		NumEntries++
		_, err := Consumer.ReadEntry()
		if err != nil {
			return NumEntries
		}
	}
}

func (fr *FileRepository) PingConnect() error {
	return nil
}
