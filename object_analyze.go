package sibyl2

import (
	"fmt"

	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type SymbolWithPath struct {
	*extractor.Symbol `bson:",inline"`
	Path              string `json:"path"`
}

// FunctionWithPath
// original symbol and function do not have a path
// because they maybe not come from a real file
type FunctionWithPath struct {
	*extractor.Function `bson:",inline"`
	Path                string `json:"path" bson:"path"`
}

func (fwp *FunctionWithPath) GetDescWithPath() string {
	return fmt.Sprintf("%s%s%s", fwp.Path, object.DescSplit, fwp.Function.GetDesc())
}

func WrapFuncWithPath(f *extractor.Function, p string) *FunctionWithPath {
	return &FunctionWithPath{
		Function: f,
		Path:     p,
	}
}

type FuncTag = string

type FunctionWithTag struct {
	*FunctionWithPath `bson:",inline"`
	Tags              []FuncTag `json:"tags" bson:"tags"`
}

func (fwt *FunctionWithTag) AddTag(tag FuncTag) {
	fwt.Tags = append(fwt.Tags, tag)
}

type ClazzWithPath struct {
	*extractor.Clazz `bson:",inline"`
	Path             string `json:"path"`
}
