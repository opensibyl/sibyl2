package core

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
)

type LangType string

const (
	JAVA   LangType = "JAVA"
	GOLANG LangType = "GOLANG"
)

func (langType *LangType) GetLanguage() *sitter.Language {
	switch *langType {
	case JAVA:
		return java.GetLanguage()
	case GOLANG:
		return golang.GetLanguage()
	}
	return nil
}

func (langType *LangType) GetFileSuffix() string {
	switch *langType {
	case JAVA:
		return ".java"
	case GOLANG:
		return ".go"
	}
	return ""
}
