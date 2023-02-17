package sibyl2

import (
	"encoding/json"
	"fmt"

	"github.com/dominikbraun/graph"
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
	Path string `json:"path"`
}

func (fwp *FunctionWithPath) GetDescWithPath() string {
	return fmt.Sprintf("<fwp %s %s>", fwp.Path, fwp.Function.GetDesc())
}

func WrapFuncWithPath(f *extractor.Function, p string) *FunctionWithPath {
	return &FunctionWithPath{
		Function: f,
		Path:     p,
	}
}

type FunctionContext struct {
	*FunctionWithPath
	Calls        []*FunctionWithPath `json:"calls"`
	ReverseCalls []*FunctionWithPath `json:"reverseCalls"`
}

// FunctionContextSlim instead of whole object, slim will only keep the signature
type FunctionContextSlim struct {
	*FunctionWithPath
	Calls        []string `json:"calls"`
	ReverseCalls []string `json:"reverseCalls"`
}

type ClazzWithPath struct {
	*extractor.Clazz
	Path string `json:"path"`
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

// ToMap export a very simple map without any custom structs. It will lose ptr to origin unit.
func (f *FunctionContext) ToMap() (map[string]any, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (f *FunctionContext) ToJson() ([]byte, error) {
	m, err := f.ToMap()
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (f *FunctionContext) ToSlim() *FunctionContextSlim {
	slim := &FunctionContextSlim{
		FunctionWithPath: f.FunctionWithPath,
		Calls:            make([]string, 0, len(f.Calls)),
		ReverseCalls:     make([]string, 0, len(f.ReverseCalls)),
	}
	for _, each := range f.Calls {
		slim.Calls = append(slim.Calls, each.GetSignature())
	}
	for _, each := range f.ReverseCalls {
		slim.ReverseCalls = append(slim.ReverseCalls, each.GetSignature())
	}
	return slim
}
