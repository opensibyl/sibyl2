package extractor

import (
	"sibyl2/pkg/core"
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

type SymbolSupport interface {
	IsSymbol(*core.Unit) bool
	ExtractSymbols([]*core.Unit) ([]*core.Symbol, error)
}

type FunctionSupport interface {
	IsFunction(*core.Unit) bool
	ExtractFunctions([]*core.Unit) ([]*core.Function, error)
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
