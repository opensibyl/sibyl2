package sibyl2

import (
	"fmt"

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

type FuncTag = string

type FunctionWithTag struct {
	*FunctionWithPath
	Tags []FuncTag `json:"tags"`
}

func WrapFuncWithTag(f *FunctionWithPath) *FunctionWithTag {
	return &FunctionWithTag{
		FunctionWithPath: f,
		Tags:             make([]FuncTag, 0),
	}
}

func (fwt *FunctionWithTag) AddTag(tag FuncTag) {
	fwt.Tags = append(fwt.Tags, tag)
}

type ClazzWithPath struct {
	*extractor.Clazz
	Path string `json:"path"`
}
