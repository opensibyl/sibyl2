package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type Runner struct {
}

func (r *Runner) ScanFiles(filePath string, lang LangType) ([]string, error) {
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

func (r *Runner) HandleFile(filePath string, lang LangType) ([]FileSymbol, error) {
	langSupport := lang.GetLanguage()
	files, err := r.ScanFiles(filePath, lang)
	if err != nil {
		return nil, err
	}

	parser := NewParser(langSupport)

	// why we use withCancel here:
	// tree-sitter has a special handler for cancelable
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var fileSymbols []FileSymbol
	fileSymbolsChan := make(chan []Symbol, len(files))
	for _, eachFile := range files {
		r.parseFileAsync(eachFile, parser, ctx, fileSymbolsChan)
	}

	for range files {
		eachFileSymbol := <-fileSymbolsChan
		if eachFileSymbol == nil {
			continue
		}
		curFileSymbol := FileSymbol{
			Path:     filePath,
			Language: lang,
			Symbols:  eachFileSymbol,
		}
		fileSymbols = append(fileSymbols, curFileSymbol)
	}
	return fileSymbols, nil
}

func (r *Runner) parseFileAsync(filepath string, parser *Parser, ctx context.Context, result chan []Symbol) {
	symbols, err := r.parseFile(filepath, parser, ctx)
	if err != nil {
		// ignore?
		result <- nil
	} else {
		result <- symbols
	}
}

func (r *Runner) parseFile(filePath string, parser *Parser, ctx context.Context) ([]Symbol, error) {
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
