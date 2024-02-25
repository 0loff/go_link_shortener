// Package app - this is a package for creating an application instance
package app

import (
	"net/http"

	"github.com/0loff/go_link_shortener/config"
	"github.com/0loff/go_link_shortener/internal/handler"
	"github.com/0loff/go_link_shortener/internal/repository"
	dbrepository "github.com/0loff/go_link_shortener/internal/repository/db_repository"
	filerepository "github.com/0loff/go_link_shortener/internal/repository/file_repository"
	inmemoryrepository "github.com/0loff/go_link_shortener/internal/repository/inmemory_repository"
	"github.com/0loff/go_link_shortener/internal/service"
)

// App - this is the main application structure
type App struct {
	Cfg        config.Config
	HttpServer *http.Server

	useCase service.Service
}

// NewApp - is the app instance initialization method
func NewApp() *App {
	app := &App{
		Cfg: config.NewConfigBuilder(),
	}

	app.useCase = *service.NewService(app.makeRepository(), app.Cfg.BaseURL)
	app.HttpServer = &http.Server{
		Addr:    app.Cfg.ServerAddress,
		Handler: handler.NewHandler(&app.useCase).InitRoutes(),
	}

	return app
}

func (a *App) makeRepository() repository.URLKeeper {
	if a.Cfg.DatabaseDSN != "" {
		return dbrepository.NewRepository(a.Cfg.DatabaseDSN)
	}

	if a.Cfg.StorageFile != "" {
		return filerepository.NewRepository(a.Cfg.StorageFile)
	}

	return inmemoryrepository.NewRepository()
}
