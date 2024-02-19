package main

import (
	"github.com/0loff/go_link_shortener/internal/app"
	"github.com/0loff/go_link_shortener/internal/utils"
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

	app.NewApp().Run()
}
