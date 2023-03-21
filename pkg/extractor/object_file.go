package extractor

import (
	"fmt"
	"path/filepath"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type BaseFileResult[T DataType] struct {
	Path     string        `json:"path"`
	Language core.LangType `json:"language"`
	Type     string        `json:"type"`
	Units    []T           `json:"units"`
}

func (b *BaseFileResult[T]) IsEmpty() bool {
	return len(b.Units) == 0
}

type FileResult = BaseFileResult[DataType]
type SymbolFileResult = BaseFileResult[*Symbol]
type FunctionFileResult = BaseFileResult[*Function]
type CallFileResult = BaseFileResult[*Call]
type ClazzFileResult = BaseFileResult[*Clazz]

func PathStandardize(results []*FileResult, basedir string) error {
	for _, each := range results {
		newPath, err := filepath.Rel(basedir, each.Path)
		if err != nil {
			return err
		}

		each.Path = filepath.ToSlash(newPath)
	}
	return nil
}

type SymbolWithPath struct {
	*Symbol `bson:",inline"`
	Path    string `json:"path"`
}

// FunctionWithPath
// original symbol and function do not have a path
// because they maybe not come from a real file
type FunctionWithPath struct {
	*Function `bson:",inline"`
	Path      string `json:"path" bson:"path"`
}

func (fwp *FunctionWithPath) GetDescWithPath() string {
	return fmt.Sprintf("%s%s%s", fwp.Path, object.DescSplit, fwp.Function.GetDesc())
}

func WrapFuncWithPath(f *Function, p string) *FunctionWithPath {
	return &FunctionWithPath{
		Function: f,
		Path:     p,
	}
}

type ClazzWithPath struct {
	*Clazz `bson:",inline"`
	Path   string `json:"path"`
}
