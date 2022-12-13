package sibyl2

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func TestQuery(t *testing.T) {
	fileResult, err := ExtractFunction("./extract.go", DefaultConfig())
	if err != nil {
		panic(err)
	}
	for _, each := range fileResult {
		for _, eachUnit := range each.Units {
			core.Log.Debugf("unit: %v", eachUnit.GetDesc())
		}
	}

	symbolResult, err := ExtractSymbol(".", DefaultConfig())
	if err != nil {
		panic(err)
	}

	// query by lines
	for _, eachFile := range fileResult {
		affectedUnits := QueryUnitsByLines(eachFile, 54, 55)
		for _, each := range affectedUnits {
			core.Log.Debugf("affected units: %s", each.GetDesc())
		}
	}

	// query by index
	for _, eachFile := range symbolResult {
		references := QueryUnitsByIndexNames(eachFile, "Extract")
		for _, each := range references {
			core.Log.Infof("found ref in %s %v: %s", eachFile.Path, each.GetSpan(), each.GetDesc())
		}
	}
}
