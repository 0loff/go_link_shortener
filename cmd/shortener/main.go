package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/pkg/handler"
	"github.com/0loff/go_link_shortener/pkg/repository"
	dbrepository "github.com/0loff/go_link_shortener/pkg/repository/dbRepository"
	filerepository "github.com/0loff/go_link_shortener/pkg/repository/fileRepository"
	inmemoryrepository "github.com/0loff/go_link_shortener/pkg/repository/inmemoryRepository"
	"github.com/0loff/go_link_shortener/pkg/service"
)

// Переменная роутера chi для инициализации во время запуска приложения
var Router chi.Router

func main() {
	NewConfigBuilder()

	services := service.NewService(makeRepository(&config), config.ShortURLHost)
	handlers := handler.NewHandler(services)

	Router = handlers.InitRoutes()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}

	logger.Sugar.Infoln("Host", config.Host)
	return http.ListenAndServe(config.Host, Router)
}

func makeRepository(cfg *Config) repository.URLKeeper {
	if config.DatabaseDSN != "" {
		return dbrepository.NewRepository(config.DatabaseDSN)
	}

	if config.StorageFile != "" {
		return filerepository.NewRepository(config.StorageFile)
	}

	return inmemoryrepository.NewRepository()
}
