package sibyl2

import (
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

// These functions are designed on the top of query.go
// for some higher levels usages
// Starts with `Analyze`

type SymbolWithPath struct {
	*extractor.Symbol        // nested
	Path              string `json:"path"`
}

// FunctionWithPath
// original symbol and function do not have a path
// because they maybe not come from a real file
type FunctionWithPath struct {
	*extractor.Function        // nested
	Path                string `json:"path"`
}

type FunctionWithRefLink struct {
	*FunctionWithPath
	Link []*FunctionWithPath `json:"link"`
}

func (fwr *FunctionWithRefLink) GetRefLinkRepr() string {
	ret := make([]string, 0, len(fwr.Link))
	for _, each := range fwr.Link {
		ret = append(ret, each.GetIndexName())
	}
	return strings.Join(ret, "<-")
}

type FuncGraphType = graph.Graph[string, *FunctionWithPath]

type FuncGraph struct {
	ReverseGraph FuncGraphType
	CallGraph    FuncGraphType
}

func (fg *FuncGraph) FindReferences(f *extractor.Function) []*FunctionWithRefLink {
	return fg.bfs(fg.ReverseGraph, f)
}

func (fg *FuncGraph) FindCalls(f *extractor.Function) []*FunctionWithRefLink {
	return fg.bfs(fg.CallGraph, f)
}

func (fg *FuncGraph) bfs(g FuncGraphType, f *extractor.Function) []*FunctionWithRefLink {
	selfDesc := f.GetDesc()
	var ret []*FunctionWithRefLink
	graph.BFS(g, f.GetDesc(), func(s string) bool {
		vertex, err := g.Vertex(s)
		// exclude itself
		if (err == nil) && (vertex.GetDesc() != selfDesc) {
			fwo := &FunctionWithRefLink{}
			fwo.FunctionWithPath = vertex
			path, err := graph.ShortestPath(g, selfDesc, vertex.GetDesc())
			if err != nil {
				// ignore this link
				return false
			}
			for _, each := range path {
				fwp, err := g.Vertex(each)
				if err != nil {
					return false
				}
				fwo.Link = append(fwo.Link, fwp)
			}
			ret = append(ret, fwo)
		}

		return false
	})

	return ret
}

func AnalyzeFuncGraph(funcFiles []*extractor.FunctionFileResult, symbolFiles []*extractor.SymbolFileResult) (*FuncGraph, error) {
	reverseGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())
	callGraph := graph.New((*FunctionWithPath).GetDesc, graph.Directed())

	// speed up cache
	funcFileMap := make(map[string]*extractor.FunctionFileResult, len(funcFiles))
	for _, each := range funcFiles {
		funcFileMap[each.Path] = each
	}

	for _, eachFuncFile := range funcFiles {
		for _, eachFunc := range eachFuncFile.Units {
			// multi graphs shared
			fwp := &FunctionWithPath{
				eachFunc,
				eachFuncFile.Path,
			}
			err := reverseGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}
			err = callGraph.AddVertex(fwp)
			if err != nil {
				return nil, err
			}

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
						reverseGraph.AddEdge(eachFunc.GetDesc(), eachMatchFunc.GetDesc())
						callGraph.AddEdge(eachMatchFunc.GetDesc(), eachFunc.GetDesc())
					}
				}
			}
		}
	}
	fg := &FuncGraph{
		ReverseGraph: reverseGraph,
		CallGraph:    callGraph,
	}
	return fg, nil
}
