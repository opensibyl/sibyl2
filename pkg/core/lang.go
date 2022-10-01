package core

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/python"
)

type LangType string

const (
	LangJava    LangType = "JAVA"
	LangGo      LangType = "GOLANG"
	LangPython  LangType = "PYTHON"
	LangUnknown LangType = "UNKNOWN"
)

func (langType LangType) GetValue() string {
	return string(langType)
}

func LangTypeValueOf(raw string) LangType {
	switch raw {
	case LangJava.GetValue():
		return LangJava
	case LangGo.GetValue():
		return LangGo
	case LangPython.GetValue():
		return LangPython
	default:
		return LangUnknown
	}
}

func (langType LangType) GetParser() *Parser {
	return NewParser(langType)
}

func (langType LangType) GetLanguage() *sitter.Language {
	switch langType {
	case LangJava:
		return java.GetLanguage()
	case LangGo:
		return golang.GetLanguage()
	case LangPython:
		return python.GetLanguage()
	}
	return nil
}

func (langType LangType) GetFileSuffix() string {
	switch langType {
	case LangJava:
		return ".java"
	case LangGo:
		return ".go"
	case LangPython:
		return ".py"
	}
	return ""
}
