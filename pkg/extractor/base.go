package extractor

import (
	"sibyl2/pkg/model"
)

/*
Extractor

in tree-sitter, a specific language is ruled by grammar.js.
https://github.com/tree-sitter/tree-sitter-java/blob/master/grammar.js

unlike other projects, we will only keep the necessary parts, not the whole grammar rule
*/
type Extractor interface {
	GetLang() model.LangType
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
	IsSymbol(*model.Unit) bool
	ExtractSymbols([]*model.Unit) ([]*model.Symbol, error)
}

type FunctionSupport interface {
	IsFunction(*model.Unit) bool
	ExtractFunctions([]*model.Unit) ([]*model.Function, error)
	ExtractFunction(*model.Unit) (*model.Function, error)
}

type CallSupport interface {
	IsCall(unit *model.Unit) bool
	ExtractCalls([]*model.Unit) ([]*model.Call, error)
}

func GetExtractor(lang model.LangType) Extractor {
	switch lang {
	case model.LangJava:
		return &JavaExtractor{}
	case model.LangGo:
		return &GolangExtractor{}
	}
	return nil
}
