package extractor

import (
	"sibyl2/pkg/core"
	"sibyl2/pkg/model"
	"testing"
)

var goCode = `
type Parser struct {
	engine *sitter.Parser
}

func NormalFunc(lang *sitter.Language) string {
	return "hello"
}

func (*Parser) NormalMethod(lang *sitter.Language) string {
	return "hi"
}
`

func TestGolangExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(model.LangGo)
	units, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(model.LangGo)
	_, err = extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
}
