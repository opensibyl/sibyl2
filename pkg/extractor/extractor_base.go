package extractor

import (
	"github.com/opensibyl/sibyl2/pkg/core"
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
}

type ExtractType = string

// these extractors are independent with each other
const (
	TypeExtractFunction ExtractType = "func"
	TypeExtractSymbol   ExtractType = "symbol"
	TypeExtractCall     ExtractType = "call"
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

type CallSupport interface {
	IsCall(unit *core.Unit) bool
	ExtractCalls([]*core.Unit) ([]*Call, error)
}

func GetExtractor(lang core.LangType) Extractor {
	switch lang {
	case core.LangJava:
		return &JavaExtractor{}
	case core.LangGo:
		return &GolangExtractor{}
	case core.LangPython:
		return &PythonExtractor{}
	}
	return nil
}
