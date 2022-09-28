package extractor

import "sibyl2/pkg/core"

type JavaExtractor struct {
}

func (extractor *JavaExtractor) IsSymbol(unit core.Unit) bool {
	// todo
	return true
}

func (extractor *JavaExtractor) ExtractSymbols(unit []core.Unit) []core.Symbol {
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
