// Пакет для логирования в приложении.
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Конструктор инициализации логгера Zap
var Log *zap.Logger = zap.NewNop()

// Упрощенный метод вызова ллоггера Zap с помощью синтаксического сахара
var Sugar *zap.SugaredLogger

// Инициализация механизма логирования в приложении
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	Sugar = Log.Sugar()

	return nil
}

// Middleware обработчик для логирования информации из запроса
func RequestLogger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
}
