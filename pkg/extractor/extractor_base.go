package extractor

import (
	"sync"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/golang"
	"github.com/opensibyl/sibyl2/pkg/extractor/java"
	"github.com/opensibyl/sibyl2/pkg/extractor/javascript"
	"github.com/opensibyl/sibyl2/pkg/extractor/kotlin"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
	"github.com/opensibyl/sibyl2/pkg/extractor/python"
)

/*
Extractor

in tree-sitter, a specific language is ruled by grammar.js.
https://github.com/tree-sitter/tree-sitter-java/blob/master/grammar.js

unlike other projects, we will only keep the necessary parts, not the whole grammar rule
*/
type Extractor interface {
	GetLang() core.LangType
	SymbolSupport
	FunctionSupport
	CallSupport
	ClassSupport
}

type ExtractType = string

var (
	extractorMu          sync.RWMutex
	additionalExtractors = make(map[core.LangType]Extractor)
)

// these extractors are independent with each other
const (
	TypeExtractFunction ExtractType = "func"
	TypeExtractSymbol   ExtractType = "symbol"
	TypeExtractCall     ExtractType = "call"
	TypeExtractClazz    ExtractType = "class"
)

type SymbolSupport interface {
	IsSymbol(*core.Unit) bool
	ExtractSymbols([]*core.Unit) ([]*Symbol, error)
}

type FunctionSupport interface {
	IsFunction(*core.Unit) bool
	ExtractFunctions([]*core.Unit) ([]*Function, error)
	ExtractFunction(*core.Unit) (*Function, error)
}

type ClassSupport interface {
	IsClass(*core.Unit) bool
	ExtractClasses([]*core.Unit) ([]*Clazz, error)
	ExtractClass(*core.Unit) (*Clazz, error)
}

type CallSupport interface {
	IsCall(unit *core.Unit) bool
	ExtractCalls([]*core.Unit) ([]*Call, error)
}

func GetExtractor(lang core.LangType) Extractor {
	switch lang {
	case core.LangJava:
		return &java.Extractor{}
	case core.LangGo:
		return &golang.Extractor{}
	case core.LangPython:
		return &python.Extractor{}
	case core.LangKotlin:
		return &kotlin.Extractor{}
	case core.LangJavaScript:
		return &javascript.Extractor{}
	}
	if e, ok := additionalExtractors[lang]; ok {
		return e
	}
	return nil
}

func RegisterExtractor(langType core.LangType, extractor Extractor) {
	extractorMu.Lock()
	defer extractorMu.Unlock()
	if extractor == nil {
		panic("extractor is nil")
	}
	if _, dup := additionalExtractors[langType]; dup {
		panic("Register called twice for lang " + langType)
	}
	additionalExtractors[langType] = extractor
}

type Function = object.Function
type Symbol = object.Symbol
type Call = object.Call
type Clazz = object.Clazz
