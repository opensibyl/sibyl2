package extractor

import (
	"sibyl2/pkg/core"
	"strings"
)

type JavaExtractor struct {
}

func (extractor *JavaExtractor) GetLang() core.LangType {
	return core.JAVA
}

func (extractor *JavaExtractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractSymbols(units []*core.Unit) ([]*core.Symbol, error) {
	var ret []*core.Symbol
	for _, eachUnit := range units {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &core.Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
}

func (extractor *JavaExtractor) IsFunction(unit *core.Unit) bool {
	if unit.Kind == "method_declaration" {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractFunctions(units []*core.Unit) ([]*core.Function, error) {
	var ret []*core.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}

		eachFunc := &core.Function{
			Name:       eachUnit.Content,
			Parameters: nil,
			Returns:    nil,
			Span:       eachUnit.Span,
		}
		ret = append(ret, eachFunc)
	}
	return ret, nil
}
