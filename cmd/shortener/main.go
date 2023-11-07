package main

import (
	"go_link_shortener/internal/logger"
	"go_link_shortener/pkg/handler"
	"go_link_shortener/pkg/repository"
	"go_link_shortener/pkg/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

var Router chi.Router

func main() {
	NewConfigBuilder()

	conn, err := pgx.ParseConfig(config.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	db, err := repository.NewPostgresDB(conn.ConnString())
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, config.ShortURLHost, config.StorageFile)
	services.SetTestShortURL()
	services.RecoverFromFile()
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
