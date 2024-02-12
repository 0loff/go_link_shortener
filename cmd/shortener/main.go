package main

import (
	"fmt"
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

// Version information variables. Initialized durign the build process.
// For example, use next command for build app shortener
//
// go build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildCommit=$(git rev-parse HEAD~1)'"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Router Переменная роутера chi для инициализации во время запуска приложения
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

	printBuildTag("version", buildVersion)
	printBuildTag("date", buildDate)
	printBuildTag("commit", buildCommit)

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

func printBuildTag(tagName, tagValue string) {
	fmt.Printf("Build %s: %s\n", tagName, tagValue)
}
