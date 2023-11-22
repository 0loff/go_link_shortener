package filerepository

import (
	"context"
	filehandler "go_link_shortener/internal/fileHandler"
	"go_link_shortener/internal/models"
	"go_link_shortener/pkg/repository"
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

func (fr *FileRepository) FindByID(ctx context.Context, id string) string {

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

func (fr *FileRepository) FindByLink(ctx context.Context, link string) string {
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

func (fr *FileRepository) FindByUser(ctx context.Context, uid string) []models.URLEntry {
	var URLEntries []models.URLEntry

	Consumer, err := filehandler.NewConsumer(fr.StorageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer Consumer.Close()

	for {
		entry, err := Consumer.ReadEntry()
		if err != nil {
			return URLEntries
		}

		if entry.UserID == uid {
			URLEntries = append(URLEntries, models.URLEntry{
				ShortURL:    entry.ShortURL,
				OriginalURL: entry.OriginalURL,
			})
		}
	}
}

func (fr *FileRepository) SetShortURL(ctx context.Context, uid, shortURL, originURL string) (string, error) {

	if fr.FindByLink(ctx, originURL) != "" {
		return shortURL, repository.ErrConflict
	}

	newEntry := filehandler.Entry{
		ID:          fr.GetNumberOfEntries(ctx),
		UserID:      uid,
		ShortURL:    shortURL,
		OriginalURL: originURL,
	}

	fr.WriteToFile(newEntry)
	return shortURL, nil
}

func (fr *FileRepository) BatchInsertShortURLS(ctx context.Context, uid string, urls []models.URLEntry) error {
	for _, u := range urls {
		fr.WriteToFile(filehandler.Entry{
			ID:          fr.GetNumberOfEntries(ctx),
			UserID:      uid,
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

func (fr *FileRepository) GetNumberOfEntries(ctx context.Context) int {
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

func (fr *FileRepository) PingConnect(ctx context.Context) error {
	return nil
}
