package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/0loff/go_link_shortener/internal/app"
	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/tls"
	"github.com/0loff/go_link_shortener/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func main() {

	utils.PrintBuildTag("version", buildVersion)
	utils.PrintBuildTag("date", buildDate)
	utils.PrintBuildTag("commit", buildCommit)

	Run(app.NewApp())
}

// Run - this is the App server startup method
func Run(a *app.App) {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		if err := logger.Initialize(a.Cfg.LogLevel); err != nil {
			log.Fatal(err)
		}

		logger.Sugar.Infoln("Host", a.Cfg.ServerAddress)

		if a.Cfg.EnableHTTPS {
			const (
				cert = "cert.pem"
				key  = "key.pem"
			)

			err := tls.TLSCertCreate(cert, key)
			if err != nil {
				log.Fatal(err)
			}

			return a.HttpServer.ListenAndServeTLS(cert, key)
		}
		return a.HttpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return a.HttpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Log.Error("The reason for stopping the server is: ", zap.Error(err))
	}
}
