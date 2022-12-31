package extractor

import (
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
)

func (extractor *PythonExtractor) IsSymbol(unit *core.Unit) bool {
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *PythonExtractor) ExtractSymbols(units []*core.Unit) ([]*Symbol, error) {
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
