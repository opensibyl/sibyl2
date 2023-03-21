package sibyl2

import (
	"encoding/json"

	"github.com/dominikbraun/graph"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionContext struct {
	*extractor.FunctionWithPath
	Calls        []*extractor.FunctionWithPath `json:"calls" bson:"calls"`
	ReverseCalls []*extractor.FunctionWithPath `json:"reverseCalls" bson:"reverseCalls"`
}

func (f *FunctionContext) ToGraph() *FuncGraphType {
	markSelf := graph.VertexAttribute("fillcolor", "red")
	markDirect := graph.VertexAttribute("fillcolor", "yellow")
	markFill := graph.VertexAttribute("style", "filled")

	ctxGraph := graph.New((*extractor.FunctionWithPath).GetIndexName, graph.Directed())
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
