package sibyl2

import (
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

// These functions are designed on the top of query.go
// for some higher levels usages
// Starts with `Analyze`

func AnalyzeFuncGraph(funcFiles []*extractor.FunctionFileResult, symbolFiles []*extractor.SymbolFileResult) (*FuncGraph, error) {
	reverseCallGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())
	callGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())

	// speed up cache
	funcFileMap := make(map[string]*extractor.FunctionFileResult, len(funcFiles))
	for _, each := range funcFiles {
		funcFileMap[each.Path] = each
	}

	// fill graph with vertex
	for _, eachFuncFile := range funcFiles {
		for _, eachFunc := range eachFuncFile.Units {
			// multi graphs shared
			fwp := &FunctionWithPath{
				eachFunc,
				eachFuncFile.Path,
			}
			err := reverseCallGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}
			err = callGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}
		}
	}

	// build relationship
	for _, eachFuncFile := range funcFiles {
		for _, eachFunc := range eachFuncFile.Units {
			// find all the refs
			var refs []*SymbolWithPath
			for _, eachSymbolFile := range symbolFiles {
				symbols := QueryUnitsByIndexNames(eachSymbolFile, eachFunc.GetIndexName())
				for _, eachSymbol := range symbols {
					refs = append(refs, &SymbolWithPath{
						Symbol: eachSymbol,
						Path:   eachSymbolFile.Path,
					})
				}
			}
			// match any functions?
			for _, eachSymbol := range refs {
				if functions, ok := funcFileMap[eachSymbol.Path]; ok {
					matched := QueryUnitsByLines(functions, eachSymbol.Span.Lines()...)
					for _, eachMatchFunc := range matched {
						// eachFunc referenced by eachMatchFunc
						reverseCallGraph.AddEdge(eachFunc.GetDesc(), eachMatchFunc.GetDesc())
						callGraph.AddEdge(eachMatchFunc.GetDesc(), eachFunc.GetDesc())
					}
				}
			}
		}
	}
	fg := &FuncGraph{
		ReverseCallGraph: reverseCallGraph,
		CallGraph:        callGraph,
	}
	return fg, nil
}
