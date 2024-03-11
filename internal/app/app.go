// Package app - this is a package for creating an application instance
package app

import (
	"net/http"

	"github.com/0loff/go_link_shortener/config"
	grpcHandler "github.com/0loff/go_link_shortener/internal/handler/grpc"
	httpHandler "github.com/0loff/go_link_shortener/internal/handler/http"
	"github.com/0loff/go_link_shortener/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	dbrepository "github.com/0loff/go_link_shortener/internal/repository/db_repository"
	filerepository "github.com/0loff/go_link_shortener/internal/repository/file_repository"
	inmemoryrepository "github.com/0loff/go_link_shortener/internal/repository/inmemory_repository"
	"github.com/0loff/go_link_shortener/internal/service"

	pb "github.com/0loff/go_link_shortener/proto"

	"github.com/0loff/go_link_shortener/internal/interceptors"
)

// App - this is the main application structure
type App struct {
	Cfg        config.Config
	HttpServer *http.Server
	GrpcServer *grpc.Server

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
		Handler: httpHandler.NewHandler(&app.useCase, app.Cfg.TrustedSubnet).InitRoutes(),
	}
	app.GrpcServer = app.makeGRPCInstance()

	return app
}

func (a *App) makeGRPCInstance() *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor))

	reflection.Register(s)
	pb.RegisterShrotenerServer(s, grpcHandler.NewHandler(&a.useCase, a.Cfg.TrustedSubnet))

	return s
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
