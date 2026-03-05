package analyser

import (
	"testing"

	"github.com/Shyyw1e/selectel-linter-test/internal/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyser(t *testing.T) {
	cfg := config.Default()
	analysistest.Run(t, analysistest.TestData(), NewAnalyser(cfg), "a")
}
