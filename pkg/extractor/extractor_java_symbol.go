package extractor

import (
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func (extractor *JavaExtractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *JavaExtractor) ExtractSymbols(units []*core.Unit) ([]*Symbol, error) {
	var ret []*Symbol
	for _, eachUnit := range units {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
			unit:      eachUnit,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
}
