package sibyl2

import (
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"testing"
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
