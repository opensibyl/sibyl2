package extractor

import (
	"sibyl2/pkg/core"
)

/*
Extractor

in tree-sitter, a specific language is ruled by grammar.js.
https://github.com/tree-sitter/tree-sitter-java/tree/master/src

for example, method_declaration

grammar: AST, for lexer.
node-types: Node desc, and static type system.
*/
type Extractor interface {
	GetLang() core.LangType
	IsSymbol(*core.Unit) bool
	ExtractSymbols([]*core.Unit) ([]*core.Symbol, error)
	IsFunction(*core.Unit) bool
	ExtractFunctions([]*core.Unit) ([]*core.Function, error)
}

func GetExtractor(lang core.LangType) Extractor {
	switch lang {
	case core.JAVA:
		return &JavaExtractor{}
	case core.GOLANG:
		return &GolangExtractor{}
	}
	return nil
}
