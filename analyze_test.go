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

	targetFuncName := "unit2Call"
	target := QueryUnitsByIndexNamesInFiles(functions, targetFuncName)
	if len(target) == 0 {
		panic(errors.New("func not found"))
	}

	references := g.FindReferences(target[0])
	calls := g.FindCalls(target[0])
	core.Log.Debugf("search func %s", targetFuncName)
	for _, each := range references {
		core.Log.Debugf("found ref %s in %s, link: %s", each.GetIndexName(), each.Path, each.GetRefLinkRepr())
	}
	for _, each := range calls {
		core.Log.Debugf("found call %s in %s, link: %s", each.GetIndexName(), each.Path, each.GetRefLinkRepr())
	}
}
