package core

// Extractor kind of filter, actually
type Extractor interface {
	Extract([]Symbol) []Symbol
}

func GetExtractor(lang *LangType) Extractor {
	switch *lang {
	case JAVA:
		return &JavaExtractor{}
	}
	return &EmptyExtractor{}
}

type JavaExtractor struct {
}

func (extractor *JavaExtractor) Extract(symbols []Symbol) []Symbol {
	// todo
	return symbols
}

type EmptyExtractor struct {
}

func (extractor *EmptyExtractor) Extract(symbols []Symbol) []Symbol {
	return symbols
}
