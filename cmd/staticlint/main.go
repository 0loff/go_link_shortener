// Пакет кастомного анализатора для проверки вызова метод os.Exit из функции main приложения
// сокращателя ссылок.
//
// Для вызова аналзатор, необходимо собрать пакет staticlint (cd cmd/staticlint/)
//
// cd cmd/staticlint/ && go build
//
// Вызвать проверку анализатором командой
//
// ./staticlint ../../...

package main

import (
	osexitlinter "github.com/0loff/go_link_shortener/pkg/osexit_linter"
	"github.com/fatih/errwrap/errwrap"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var mychecks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		if v.Analyzer.Name == "ST1019" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	mychecks = append(mychecks,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		// A couple of public linters
		errwrap.Analyzer,
		ineffassign.Analyzer,
		osexitlinter.OSExitAnalyzer,
	)

	multichecker.Main(
		mychecks...,
	)
}
