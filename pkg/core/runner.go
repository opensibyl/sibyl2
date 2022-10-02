package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type Runner struct {
}

func (r *Runner) File2Units(filePath string, lang LangType) ([]*FileUnit, error) {
	files, err := r.scanFiles(filePath, lang)
	if err != nil {
		return nil, err
	}

	parser := NewParser(lang)

	// why we use withCancel here:
	// tree-sitter has a special handler for cancelable
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var fileUnits []*FileUnit
	fileUnitsChan := make(chan *FileUnit, len(files))
	for _, eachFile := range files {
		r.parseFileAsync(eachFile, parser, ctx, fileUnitsChan)
	}

	for range files {
		eachFileUnit := <-fileUnitsChan
		if eachFileUnit == nil {
			continue
		}
		eachFileUnit.Language = lang
		fileUnits = append(fileUnits, eachFileUnit)
	}

	return fileUnits, nil
}

func (r *Runner) scanFiles(filePath string, lang LangType) ([]string, error) {
	var files []string
	fileSuffix := lang.GetFileSuffix()
	handleFunc := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, fileSuffix) {
			files = append(files, path)
		}
		return nil
	}
	err := filepath.Walk(filePath, handleFunc)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *Runner) parseFileAsync(filepath string, parser *Parser, ctx context.Context, result chan *FileUnit) {
	units, err := r.parseFile(filepath, parser, ctx)
	if err != nil {
		// ignore?
		result <- nil
	} else {
		ret := &FileUnit{
			Path:     filepath,
			Language: "",
			Units:    units,
		}
		result <- ret
	}
}

func (r *Runner) parseFile(filePath string, parser *Parser, ctx context.Context) ([]*Unit, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	parsed, err := parser.ParseCtx(content, ctx)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
