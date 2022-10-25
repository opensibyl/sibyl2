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
		core.Log.Infof("path: %v, %v", each.Path, each.Units)
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
		core.Log.Infof("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractFile(t *testing.T) {
	fileResult, err := Extract("./api.go", &ExtractConfig{
		LangType:    core.LangGo,
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Infof("path: %v, %v", each.Path, each.Units)
	}
}

func TestExtractFileWithGuess(t *testing.T) {
	fileResult, err := Extract("./api.go", &ExtractConfig{
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		for _, eachUnit := range each.Units {
			core.Log.Infof("unit: %v", eachUnit.GetDesc())
		}
	}

	// query api
	for _, eachFile := range fileResult {
		affectedUnits := QueryAffectedUnitsByLine(eachFile, 54, 55)
		for _, each := range affectedUnits {
			core.Log.Infof("affected units: %s", each.GetDesc())
		}
	}
}
