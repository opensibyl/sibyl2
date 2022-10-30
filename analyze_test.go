package sibyl2

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
	"testing"
)

func TestAnalyzeFuncGraph(t *testing.T) {
	symbols, _ := ExtractSymbol(".", DefaultConfig())
	functions, _ := ExtractFunction(".", DefaultConfig())
	g, err := AnalyzeFuncGraph(functions, symbols)
	if err != nil {
		panic(err)
	}

	target := QueryUnitsByIndexNamesInFiles(functions, "NewParser")
	if len(target) == 0 {
		panic(errors.New("func not found"))
	}

	related := g.FindRelated(target[0])
	for _, each := range related {
		core.Log.Infof("found ref %s in %s", each.GetIndexName(), each.Path)
	}
}
