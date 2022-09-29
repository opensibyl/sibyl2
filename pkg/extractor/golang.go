package extractor

import (
	"sibyl2/pkg/core"
	"strings"
)

type GolangExtractor struct {
}

func (extractor *GolangExtractor) GetLang() core.LangType {
	return core.GOLANG
}

func (extractor *GolangExtractor) IsSymbol(unit core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractSymbols(unit []core.Unit) []core.Symbol {
	var ret []core.Symbol
	for _, eachUnit := range unit {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := core.Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
			// todo
			NodeType:   "",
			SyntaxType: "",
		}
		ret = append(ret, symbol)
	}
	return ret
}

func (extractor *GolangExtractor) IsFunction(unit core.Unit) bool {
	if unit.Kind == "function_declaration" {
		return true
	}
	return false
}

func (extractor *GolangExtractor) ExtractFunctions(units []core.Unit) []core.Function {
	var ret []core.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc := core.Function{
			Name:       eachUnit.Content,
			Parameters: nil,
			Returns:    nil,
			Span:       eachUnit.Span,
		}
		ret = append(ret, eachFunc)
	}
	return ret
}
