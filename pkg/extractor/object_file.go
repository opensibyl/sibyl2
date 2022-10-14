package extractor

import (
	"github.com/williamfzc/sibyl2/pkg/core"
	"path/filepath"
)

type FileResult struct {
	Path     string        `json:"path"`
	Language core.LangType `json:"language"`
	Type     string        `json:"type"`
	Units    []DataType    `json:"units"`
}

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
