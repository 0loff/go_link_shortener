package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/0loff/go_link_shortener/internal/handler"
	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/repository"
	dbrepository "github.com/0loff/go_link_shortener/internal/repository/db_repository"
	filerepository "github.com/0loff/go_link_shortener/internal/repository/file_repository"
	inmemoryrepository "github.com/0loff/go_link_shortener/internal/repository/inmemory_repository"
	"github.com/0loff/go_link_shortener/internal/service"
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
