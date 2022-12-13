package sibyl2

import (
	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type SymbolWithPath struct {
	*extractor.Symbol
	Path string `json:"path"`
}

// FunctionWithPath
// original symbol and function do not have a path
// because they maybe not come from a real file
type FunctionWithPath struct {
	*extractor.Function
	Path     string        `json:"path"`
	Language core.LangType `json:"language"`
}

type FunctionContext struct {
	*FunctionWithPath
	Calls        []*FunctionWithPath `json:"calls"`
	ReverseCalls []*FunctionWithPath `json:"reverseCalls"`
}

func (f *FunctionContext) ToGraph() *FuncGraphType {
	markSelf := graph.VertexAttribute("fillcolor", "red")
	markDirect := graph.VertexAttribute("fillcolor", "yellow")
	markFill := graph.VertexAttribute("style", "filled")

	ctxGraph := graph.New((*FunctionWithPath).GetIndexName, graph.Directed())
	_ = ctxGraph.AddVertex(f.FunctionWithPath, markFill, markSelf)
	for _, each := range f.Calls {
		// bind itself
		_ = ctxGraph.AddVertex(each, markFill, markDirect)
		_ = ctxGraph.AddEdge(f.GetIndexName(), each.GetIndexName())
	}
	for _, each := range f.ReverseCalls {
		// bind itself
		_ = ctxGraph.AddVertex(each, markFill, markDirect)
		_ = ctxGraph.AddEdge(each.GetIndexName(), f.GetIndexName())
	}
	return WrapFuncGraph(ctxGraph)
}
