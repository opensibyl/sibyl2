package extractor

import (
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"
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

func Abcd[T DataType](result *BaseFileResult[T]) []T {
	return nil
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
