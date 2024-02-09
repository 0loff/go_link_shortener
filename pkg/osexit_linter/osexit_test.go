package osexitlinter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOSExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), OSExitAnalyzer, "./...")
}
