// Package app - this is a package for creating an application instance
package app

import (
	"log"
	"net/http"

	"github.com/0loff/go_link_shortener/config"
	"github.com/0loff/go_link_shortener/internal/handler"
	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/repository"
	dbrepository "github.com/0loff/go_link_shortener/internal/repository/db_repository"
	filerepository "github.com/0loff/go_link_shortener/internal/repository/file_repository"
	inmemoryrepository "github.com/0loff/go_link_shortener/internal/repository/inmemory_repository"
	"github.com/0loff/go_link_shortener/internal/service"
	"github.com/0loff/go_link_shortener/internal/utils"
)

type app struct {
	cfg        config.Config
	httpServer *http.Server

	useCase service.Service
}

// NewApp - is the app instance initialization method
func NewApp() *app {
	app := &app{
		cfg: config.NewConfigBuilder(),
	}

	app.useCase = *service.NewService(app.makeRepository(), app.cfg.ShortURLHost)
	app.httpServer = &http.Server{
		Addr:    app.cfg.Host,
		Handler: handler.NewHandler(&app.useCase).InitRoutes(),
	}

	return app
}

// Run - this method for run app server
func (a *app) Run() error {
	if err := logger.Initialize(a.cfg.LogLevel); err != nil {
		return err
	}

	logger.Sugar.Infoln("Host", a.cfg.Host)

	if a.cfg.EnableHTTPS {
		const (
			cert = "cert.pem"
			key  = "key.pem"
		)

		err := utils.TLSCertCreate(cert, key)
		if err != nil {
			log.Fatal(err)
		}

		return a.httpServer.ListenAndServeTLS(cert, key)
	}

	return a.httpServer.ListenAndServe()
}

func (a *app) makeRepository() repository.URLKeeper {
	if a.cfg.DatabaseDSN != "" {
		return dbrepository.NewRepository(a.cfg.DatabaseDSN)
	}

	if a.cfg.StorageFile != "" {
		return filerepository.NewRepository(a.cfg.StorageFile)
	}

	return inmemoryrepository.NewRepository()
}
