package core

import (
	"context"
	"os"
	"path/filepath"
)

/*
Runner

binding to file system
*/
type Runner struct {
}

func (r *Runner) File2Units(path string, lang LangType, fileFilter func(string) bool) ([]*FileUnit, error) {
	files, err := r.scanFiles(path, lang, fileFilter)
	if err != nil {
		return nil, err
	}
	Log.Infof("valid file count: %d", len(files))

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

	// wait for all the tasks finished
	for range files {
		eachFileUnit := <-fileUnitsChan
		if eachFileUnit == nil {
			continue
		}
		eachFileUnit.Language = lang
		fileUnits = append(fileUnits, eachFileUnit)
		Log.Debugf("collect units: %d from file: %s", len(eachFileUnit.Units), eachFileUnit.Path)
	}

	return fileUnits, nil
}

func (r *Runner) scanFiles(filePath string, lang LangType, fileFilter func(string) bool) ([]string, error) {
	var files []string

	handleFunc := func(path string, info os.FileInfo, err error) error {
		if !lang.MatchName(path) {
			return nil
		}
		if fileFilter != nil {
			if !fileFilter(path) {
				return nil
			}
		}
		files = append(files, path)
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
		Log.Errorf("error when parse file %s, err: %v", filepath, err)
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

func (r *Runner) GuessLangFromDir(dir string, fileFilter func(string) bool) (LangType, error) {
	countMap := make(map[LangType]int, len(SupportedLangs))
	for _, each := range SupportedLangs {
		countMap[each] = 0
	}

	handleFunc := func(path string, info os.FileInfo, err error) error {
		if fileFilter != nil {
			if !fileFilter(path) {
				return nil
			}
		}
		for k := range countMap {
			if k.MatchName(path) {
				countMap[k]++
			}
		}
		return nil
	}
	err := filepath.Walk(dir, handleFunc)
	if err != nil {
		return "", err
	}

	ret := LangUnknown
	max := 0
	for k := range countMap {
		if countMap[k] > max {
			ret = k
		}
	}
	return ret, nil
}
