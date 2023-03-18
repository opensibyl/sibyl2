package extractor

import (
	"path/filepath"

	"github.com/opensibyl/sibyl2/pkg/core"
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
