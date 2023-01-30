package sibyl2

import (
	"sync"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type AdjacencyMapType = map[string]map[string]graph.Edge[string]

type FuncGraphType struct {
	graph.Graph[string, *FunctionWithPath]
	adjMapCache *AdjacencyMapType
	l           *sync.Mutex
}

func WrapFuncGraph(g graph.Graph[string, *FunctionWithPath]) *FuncGraphType {
	return &FuncGraphType{
		g,
		nil,
		&sync.Mutex{},
	}
}

/*
FuncGraph

It is not a serious `call` graph.
It based on references not real calls.

Why we used it:
- We can still use something like `method_invocation`
- But we mainly use it to evaluate the influence of a method
- In many languages, scope of `invocation` is too small
- For example, use `function` as a parameter.
*/
type FuncGraph struct {
	ReverseCallGraph *FuncGraphType
	CallGraph        *FuncGraphType
}

func (fg *FuncGraph) FindReverseCalls(f *extractor.Function) []*FunctionWithPath {
	return fg.bfs(fg.ReverseCallGraph, f)
}

func (fg *FuncGraph) FindCalls(f *extractor.Function) []*FunctionWithPath {
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

func (fg *FuncGraph) bfs(g *FuncGraphType, f *extractor.Function) []*FunctionWithPath {
	selfDesc := f.GetDesc()
	ret := make([]*FunctionWithPath, 0)

	// if there is an edge (a, b),
	// b is an adjacency of a.
	// but a isn't an adjacency of b.
	adjacencyMap, err := g.GetAdjacencyMap()
	if err != nil {
		return ret
	}

	// calc the shortest path can be slow in large scale graph
	// these heavy calculations should be done outside this lib
	m := (*adjacencyMap)[selfDesc]
	for k := range m {
		vertex, err := g.Vertex(k)
		if err != nil {
			core.Log.Warnf("invalid %s vertex found: %v", k, err)
			continue
		}
		ret = append(ret, vertex)
	}

	return ret
}

func (fgt *FuncGraphType) GetAdjacencyMap() (*AdjacencyMapType, error) {
	fgt.l.Lock()
	defer fgt.l.Unlock()

	cache := fgt.adjMapCache
	if cache != nil {
		return cache, nil
	}

	// rebuild cache
	m, err := fgt.AdjacencyMap()
	if err != nil {
		core.Log.Warnf("failed to get adjacency map: %v", err)
		return nil, err
	}
	fgt.adjMapCache = &m
	return &m, nil
}
