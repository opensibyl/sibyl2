package core

// Extractor kind of filter, actually
type Extractor interface {
	IsSymbol(Unit) bool
	ExtractSymbols([]Unit) []Symbol
}

type AbstractExtractor struct {
}

func (extractor *AbstractExtractor) IsSymbol(unit Unit) bool {
	panic("ABSTRACT_METHOD")
}

func (extractor *AbstractExtractor) ExtractSymbols(unit []Unit) []Symbol {
	var ret []Symbol
	for _, eachUnit := range unit {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := Symbol{
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

func GetExtractor(lang LangType) Extractor {
	switch lang {
	case JAVA:
		return &JavaExtractor{}
	case GOLANG:
		return &GolangExtractor{}
	}
	return nil
}

type JavaExtractor struct {
	*AbstractExtractor
}

func (extractor *JavaExtractor) IsSymbol(unit Unit) bool {
	// todo
	return true
}

type GolangExtractor struct {
	*AbstractExtractor
}

func (extractor *GolangExtractor) IsSymbol(unit Unit) bool {
	// todo
	return true
}
