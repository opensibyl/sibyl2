package extractor

import "sibyl2/pkg/core"

// Extractor kind of filter, actually
type Extractor interface {
	IsSymbol(core.Unit) bool
	ExtractSymbols([]core.Unit) []core.Symbol
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
