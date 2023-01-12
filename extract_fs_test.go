package sibyl2

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
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

func TestExtractClazz(t *testing.T) {
	fileResult, err := ExtractClazz("./extract.go", &ExtractConfig{
		LangType: core.LangGo,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Infof("path: %v, %v", each.Path, each.Units)
	}
}

func BenchmarkExtract(b *testing.B) {
	// with cache: 79614514 ns/op
	// no   cache: 294940375 ns/op
	for i := 0; i < b.N; i++ {
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
}
