package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// Конструктор инициализации переопределения стандартного ResponseWriter
// на ResponseWriter со сжатием gzip
func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Метод получения Header из запроса в рамках middleware обработчика,
// для последующего определения необхоимости кодирования ответа
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Метод для записи ответа
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// Метод записи Header ответа
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

// Метод закрытия потока чтения из body запроса
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// Конструктор инициализации переопределения стандартного средства чтения потока
// на средство чтения с расшифровкой сжатия формата gzip
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Чтение кодированных данных в bytes
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Закрытие потока чтения
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}
