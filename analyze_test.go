package sibyl2

import (
	"errors"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func TestAnalyzeFuncGraph(t *testing.T) {
	t.Parallel()
	symbols, _ := ExtractSymbol(".", DefaultConfig())
	functions, _ := ExtractFunction(".", DefaultConfig())
	g, err := AnalyzeFuncGraph(functions, symbols)
	if err != nil {
		panic(err)
	}

	targetFuncName := "unit2Call"
	path := "pkg/extractor/golang/extractor_golang.go"
	target := QueryUnitsByIndexNamesInFiles(functions, targetFuncName)
	if len(target) == 0 {
		panic(errors.New("func not found"))
	}

	references := g.FindReverseCalls(WrapFuncWithPath(target[0], path))
	calls := g.FindCalls(WrapFuncWithPath(target[0], path))
	core.Log.Debugf("search func %s", targetFuncName)
	for _, each := range references {
		core.Log.Debugf("found ref %s in %s", each.GetIndexName(), each.Path)
	}
	for _, each := range calls {
		core.Log.Debugf("found call %s in %s", each.GetIndexName(), each.Path)
	}
}

func TestAnalyzeFuncGraph2(t *testing.T) {
	t.Parallel()
	symbols, _ := ExtractSymbol(".", DefaultConfig())
	functions, _ := ExtractFunction(".", DefaultConfig())
	g, err := AnalyzeFuncGraph(functions, symbols)
	if err != nil {
		panic(err)
	}

	targetFuncName := "unit2Call"
	path := "pkg/extractor/golang/extractor_golang.go"
	target := QueryUnitsByIndexNamesInFiles(functions, targetFuncName)
	if len(target) == 0 {
		panic(errors.New("func not found"))
	}

	ctx := g.FindRelated(WrapFuncWithPath(target[0], path))
	for _, each := range ctx.Calls {
		core.Log.Infof("call: %s", each.Name)
	}
	for _, each := range ctx.ReverseCalls {
		core.Log.Infof("ref: %s", each.Name)
	}
}
