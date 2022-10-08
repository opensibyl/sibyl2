package pkg

import (
	"sibyl2/pkg/core"
	"sibyl2/pkg/extractor"
	"sibyl2/pkg/model"
	"testing"
)

func TestExtract(t *testing.T) {
	fileResult, err := SibylApi.Extract(".", model.LangGo, extractor.TypeExtractFunction)
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		core.Log.Infof("path: %v, %v", each.Path, each.Units)
	}
}
