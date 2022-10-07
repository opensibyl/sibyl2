package extractor

import (
	"sibyl2/pkg/core"
	"sibyl2/pkg/model"
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
}

type ExtractType = string

const (
	TypeExtractFunction ExtractType = "func"
	TypeExtractSymbol   ExtractType = "symbol"
)

type SymbolSupport interface {
	IsSymbol(*model.Unit) bool
	ExtractSymbols([]*model.Unit) ([]*model.Symbol, error)
}

type FunctionSupport interface {
	IsFunction(*model.Unit) bool
	ExtractFunctions([]*model.Unit) ([]*model.Function, error)
}

func GetExtractor(lang core.LangType) Extractor {
	switch lang {
	case core.LangJava:
		return &JavaExtractor{}
	case core.LangGo:
		return &GolangExtractor{}
	}
	return nil
}
