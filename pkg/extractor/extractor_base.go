package extractor

import (
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/golang"
	"github.com/opensibyl/sibyl2/pkg/extractor/java"
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
	}
	return nil
}

type Function = object.Function
type Symbol = object.Symbol
type Call = object.Call
type Clazz = object.Clazz

var Map2Func = object.Map2Func
