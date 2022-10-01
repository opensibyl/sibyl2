package extractor

import (
	"sibyl2/pkg/core"
	"testing"
)

var goCode = `
type Parser struct {
	engine *sitter.Parser
}

func NewParser(lang *sitter.Language) *Parser {
	engine := sitter.NewParser()
	engine.SetLanguage(lang)
	return &Parser{
		engine,
	}
}

func (p *Parser) OldParser(lang *sitter.Language) *Parser {
	engine := sitter.NewParser()
	engine.SetLanguage(lang)
	return &Parser{
		engine,
	}
} 
`

func TestGolangExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.LangGo)
	units, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.LangGo)
	_, err = extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
}
