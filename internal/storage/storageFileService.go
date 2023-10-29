package storage

import (
	"encoding/json"
	"os"
)

type (
	Entry struct {
		ID          int    `json:"uuid"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"origin_url"`
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

func (p *Producer) WriteEntry(entry *Entry) error {
	return p.encoder.Encode(&entry)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

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

func (c *Consumer) ReadEntry() (*Entry, error) {
	entry := &Entry{}
	if err := c.decoder.Decode(&entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
