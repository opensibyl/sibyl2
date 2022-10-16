package core

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"golang.org/x/exp/slices"
	"strings"
)

type LangType string

const (
	LangJava    LangType = "JAVA"
	LangGo      LangType = "GOLANG"
	LangUnknown LangType = "UNKNOWN"
)

var SupportedLangs = []LangType{
	LangJava,
	LangGo,
}

func (langType LangType) IsSupported() bool {
	return slices.Contains(SupportedLangs, langType)
}

func (langType LangType) GetValue() string {
	return string(langType)
}

func LangTypeValueOf(raw string) LangType {
	switch raw {
	case LangJava.GetValue():
		return LangJava
	case LangGo.GetValue():
		return LangGo
	default:
		return LangUnknown
	}
}

func (langType LangType) GetLanguage() *sitter.Language {
	switch langType {
	case LangJava:
		return java.GetLanguage()
	case LangGo:
		return golang.GetLanguage()
	}
	return nil
}

func (langType LangType) GetFileSuffix() string {
	switch langType {
	case LangJava:
		return ".java"
	case LangGo:
		return ".go"
	}
	return ""
}

func (langType LangType) MatchName(name string) bool {
	return strings.HasSuffix(name, langType.GetFileSuffix())
}
