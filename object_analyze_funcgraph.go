package sibyl2

import (
	"github.com/dominikbraun/graph"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

type FuncGraphType = graph.Graph[string, *FunctionWithPath]

/*
FuncGraph

It is not a serious `call` graph.
It based on references not real calls.
*/
type FuncGraph struct {
	ReverseCallGraph FuncGraphType
	CallGraph        FuncGraphType
}

func (fg *FuncGraph) FindReverseCalls(f *extractor.Function) []*FunctionWithRefLink {
	return fg.bfs(fg.ReverseCallGraph, f)
}

func (fg *FuncGraph) FindCalls(f *extractor.Function) []*FunctionWithRefLink {
	return fg.bfs(fg.CallGraph, f)
}

func (fg *FuncGraph) WrapFuncWithPath(f *extractor.Function) (*FunctionWithPath, error) {
	vertex, err := fg.CallGraph.Vertex(f.GetDesc())
	if err != nil {
		return nil, err
	}
	return vertex, nil
}

func (fg *FuncGraph) FindRelated(f *extractor.Function) *FunctionContext {
	ctx := &FunctionContext{}
	fwp, err := fg.WrapFuncWithPath(f)
	if err != nil {
		return nil
	}
	reverseCalls := fg.FindReverseCalls(f)
	calls := fg.FindCalls(f)
	ctx.FunctionWithPath = fwp
	ctx.ReverseCalls = reverseCalls
	ctx.Calls = calls
	return ctx
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
			// calc the shortest path can be slow in large scale graph
			// these heavy calculations should be done outside this lib
			fwo.Link = append(fwo.Link, vertex)
			ret = append(ret, fwo)
		}

		return false
	})

	return ret
}
