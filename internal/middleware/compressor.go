package middleware

import (
	"net/http"
	"strings"

	"github.com/0loff/go_link_shortener/pkg/compressor"
)

// Middleware обработчки запроса, позволяющий декодировать входящие данные,
// а так же кодировать исходящий ответ
func GzipCompressor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			cw := compressor.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cr, err := compressor.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	})
}
