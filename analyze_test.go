package sibyl2

import (
	"errors"
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"
)

func TestAnalyzeFuncGraph(t *testing.T) {
	symbols, _ := ExtractSymbol(".", DefaultConfig())
	functions, _ := ExtractFunction(".", DefaultConfig())
	g, err := AnalyzeFuncGraph(functions, symbols)
	if err != nil {
		panic(err)
	}

	targetFuncName := "NewParser"
	target := QueryUnitsByIndexNamesInFiles(functions, targetFuncName)
	if len(target) == 0 {
		panic(errors.New("func not found"))
	}

	related := g.FindRelated(target[0])
	core.Log.Infof("search func %s", targetFuncName)
	for _, each := range related {
		core.Log.Infof("found ref %s in %s, link: %s", each.GetIndexName(), each.Path, each.GetRefLinkRepr())
	}
}
