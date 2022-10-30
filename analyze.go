package sibyl2

import (
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

// These functions are designed on the top of query.go
// for some higher levels usages
// Starts with `Analyze`

type SymbolWithPath struct {
	*extractor.Symbol
	Path string
}

// FunctionWithPath
// original symbol and function do not have a path
// because they maybe not come from a real file
type FunctionWithPath struct {
	*extractor.Function
	Path string
}

type FuncGraph struct {
	graph.Graph[string, *FunctionWithPath]
}

func (fg *FuncGraph) FindRelated(f *extractor.Function) []*FunctionWithPath {
	var ret []*FunctionWithPath
	graph.DFS(fg.Graph, f.GetDesc(), func(s string) bool {
		vertex, err := fg.Vertex(s)
		if err == nil {
			ret = append(ret, vertex)
		}
		return false
	})
	return ret
}

func AnalyzeFuncGraph(funcFiles []*extractor.FunctionFileResult, symbolFiles []*extractor.SymbolFileResult) (*FuncGraph, error) {
	funcGraph := graph.New((*FunctionWithPath).GetDesc)

	// speed up cache
	funcFileMap := make(map[string]*extractor.FunctionFileResult, len(funcFiles))
	for _, each := range funcFiles {
		funcFileMap[each.Path] = each
	}

	for _, eachFuncFile := range funcFiles {
		for _, eachFunc := range eachFuncFile.Units {
			err := funcGraph.AddVertex(&FunctionWithPath{
				eachFunc,
				eachFuncFile.Path,
			})
			if err != nil {
				return nil, err
			}

			// find all the refs
			var refs []*SymbolWithPath
			for _, eachSymbolFile := range symbolFiles {
				symbols := QueryUnitsByIndexNames(eachSymbolFile, eachFunc.GetIndexName())
				for _, v := range symbols {
					refs = append(refs, &SymbolWithPath{
						Symbol: v,
						Path:   eachSymbolFile.Path,
					})
				}
			}
			// match any functions?
			for _, eachSymbol := range refs {
				if functions, ok := funcFileMap[eachSymbol.Path]; ok {
					matched := QueryUnitsByLines(functions, eachSymbol.Span.Lines()...)
					for _, eachMatchFunc := range matched {
						funcGraph.AddEdge(eachFunc.GetDesc(), eachMatchFunc.GetDesc())
					}
				}
			}
		}
	}
	return &FuncGraph{funcGraph}, nil
}
