package sibyl2

import (
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

func TestExtract(t *testing.T) {
	fileResult, err := Extract(".", &ExtractConfig{
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractWithGuess(t *testing.T) {
	fileResult, err := Extract(".", &ExtractConfig{
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractFunctionWithGuess(t *testing.T) {
	fileResult, err := ExtractFunction(".", &ExtractConfig{})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractSymbol(t *testing.T) {
	fileResult, err := ExtractSymbol(".", &ExtractConfig{})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %s, %v", each.Path, each.Units)
	}
}

func TestExtractFile(t *testing.T) {
	fileResult, err := Extract("./extract.go", &ExtractConfig{
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Debugf("path: %v, %v", each.Path, each.Units)
	}
}
