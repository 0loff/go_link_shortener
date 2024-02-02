package filehandler

import (
	"encoding/json"
	"io"
	"os"
)

// Структура атомарной записи для хранения в файле
type (
	Entry struct {
		ID          int    `json:"uuid"`
		UserID      string `json:"user_id"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"origin_url"`
		IsDeleted   bool   `json:"id_deleted"`
	}

	Producer struct {
		file    *os.File
		encoder *json.Encoder
	}

	Consumer struct {
		file    *os.File
		decoder *json.Decoder
	}
)

// Конструктор создания (открытия) файла для записи с возможностью преобразования
// записи в json формат
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// Создание записи в файле
func (p *Producer) WriteEntry(entry *Entry) error {
	return p.encoder.Encode(entry)
}

// Очистка файла от сохраненных записей
func (p *Producer) Trunc() {
	p.file.Seek(0, io.SeekStart)
	p.file.Truncate(0)
}

// Закрытие файла
func (p *Producer) Close() error {
	return p.file.Close()
}

// Конструктор открытия файла для чтения сохраненных записей с возможностью
// декодирования из формата json
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// Чтение записи
func (c *Consumer) ReadEntry() (*Entry, error) {
	entry := &Entry{}
	if err := c.decoder.Decode(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// Закрытие файла
func (c *Consumer) Close() error {
	return c.file.Close()
}
