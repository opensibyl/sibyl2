package core

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LangType = string

const (
	JAVA   LangType = "JAVA"
	GOLANG LangType = "GOLANG"
)

type Runner struct {
}

func (r *Runner) HandleFile(filePath string, lang LangType) ([]FileSymbol, error) {
	// lang check
	var langSupport *sitter.Language
	var fileSuffix string
	switch lang {
	case JAVA:
		langSupport = java.GetLanguage()
		fileSuffix = ".java"
	case GOLANG:
		langSupport = golang.GetLanguage()
		fileSuffix = ".go"
	}

	var files []string
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

	parser := NewParser(langSupport)
	var wg sync.WaitGroup
	wg.Add(len(files))
	ctx := context.Background()

	var ret []FileSymbol
	for _, eachFile := range files {

		content, err := os.ReadFile(eachFile)
		if err != nil {
			return nil, err
		}
		parsed, err := parser.ParseCtx(content, ctx)
		if err != nil {
			return nil, err
		}
		curFileSymbol := FileSymbol{
			Path:     eachFile,
			Language: lang,
			Symbols:  parsed,
		}
		ret = append(ret, curFileSymbol)
	}
	return ret, nil
}
