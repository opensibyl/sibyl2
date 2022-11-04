package sibyl2

import (
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
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

type FunctionWithRefLink struct {
	*FunctionWithPath
	Link []*FunctionWithPath `json:"link"` // this link will include itself
}

func (fwr *FunctionWithRefLink) GetRefLinkRepr() string {
	ret := make([]string, 0, len(fwr.Link))
	for _, each := range fwr.Link {
		ret = append(ret, each.GetIndexName())
	}
	return strings.Join(ret, "->")
}

type FunctionContext struct {
	*FunctionWithPath
	Calls        []*FunctionWithRefLink `json:"calls"`
	ReverseCalls []*FunctionWithRefLink `json:"reverseCalls"`
}

func (f *FunctionContext) ToGraph() FuncGraphType {
	markSelf := graph.VertexAttribute("fillcolor", "red")
	markDirect := graph.VertexAttribute("fillcolor", "yellow")
	markFill := graph.VertexAttribute("style", "filled")

	ctxGraph := graph.New((*FunctionWithPath).GetIndexName, graph.Directed())
	ctxGraph.AddVertex(f.FunctionWithPath, markFill, markSelf)
	for _, each := range f.Calls {
		// bind itself
		ctxGraph.AddVertex(each.FunctionWithPath, markFill, markDirect)
		ctxGraph.AddEdge(f.GetIndexName(), each.GetIndexName())
		for index := range each.Link[:len(each.Link)-1] {
			ctxGraph.AddVertex(each.Link[index])
			ctxGraph.AddEdge(each.Link[index].GetIndexName(), each.Link[index+1].GetIndexName())
		}
	}
	for _, each := range f.ReverseCalls {
		// bind itself
		ctxGraph.AddVertex(each.FunctionWithPath, markFill, markDirect)
		ctxGraph.AddEdge(each.GetIndexName(), f.GetIndexName())
		for index := range each.Link[:len(each.Link)-1] {
			ctxGraph.AddVertex(each.Link[index])
			ctxGraph.AddEdge(each.Link[index+1].GetIndexName(), each.Link[index].GetIndexName())
		}
	}
	return ctxGraph
}
