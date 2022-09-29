package core

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/python"
)

type LangType string

const (
	JAVA   LangType = "JAVA"
	GOLANG LangType = "GOLANG"
	PYTHON LangType = "PYTHON"
)

func (langType LangType) GetParser() *Parser {
	return NewParser(langType)
}

func (langType LangType) GetLanguage() *sitter.Language {
	switch langType {
	case JAVA:
		return java.GetLanguage()
	case GOLANG:
		return golang.GetLanguage()
	case PYTHON:
		return python.GetLanguage()
	}
	return nil
}

func (langType LangType) GetFileSuffix() string {
	switch langType {
	case JAVA:
		return ".java"
	case GOLANG:
		return ".go"
	case PYTHON:
		return ".py"
	}
	return ""
}
