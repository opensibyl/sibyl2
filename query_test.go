package sibyl2

import (
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"testing"
)

func TestQuery(t *testing.T) {
	fileResult, err := Extract("./extract.go", &ExtractConfig{
		ExtractType: extractor.TypeExtractFunction,
	})
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		for _, eachUnit := range each.Units {
			core.Log.Debugf("unit: %v", eachUnit.GetDesc())
		}
	}

	// query by lines
	for _, eachFile := range fileResult {
		affectedUnits := QueryUnitsByLines(eachFile, 54, 55)
		for _, each := range affectedUnits {
			core.Log.Debugf("affected units: %s", each.GetDesc())
		}
	}

	// query by index
	for _, eachFile := range fileResult {
		functions := QueryUnitsByIndexName(eachFile, "Extract")
		for _, each := range functions {
			core.Log.Infof("found func: %s", each.GetDesc())
		}
	}
}
