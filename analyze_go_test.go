package sibyl2

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	extractor2 "github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/extractor/golang"
	"github.com/stretchr/testify/assert"
)

var goCodeForAnalyze = `
package abc

type Parser struct {
	*Headless
	engine *sitter.Parser
}

func NormalFunc(lang *sitter.Language) string {
	called()
}

func called() {
	return "hello"
}
`

func TestAnalyzeGolang(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangGo)
	units, err := parser.Parse([]byte(goCodeForAnalyze))
	if err != nil {
		panic(err)
	}

	extractor := &golang.Extractor{}
	symbols, err := extractor.ExtractSymbols(units)
	functions, err := extractor.ExtractFunctions(units)
	symbolWrap := &extractor2.SymbolFileResult{}
	symbolWrap.Units = symbols
	functionWrap := &extractor2.FunctionFileResult{}
	functionWrap.Units = functions

	if err != nil {
		panic(err)
	}

	g, err := AnalyzeFuncGraph([]*extractor2.FunctionFileResult{functionWrap}, []*extractor2.SymbolFileResult{symbolWrap})
	if err != nil {
		panic(err)
	}

	ctx := g.FindRelated(functions[1])
	assert.Equal(t, ctx.Name, "called")
	assert.Empty(t, ctx.Calls)
	assert.Len(t, ctx.ReverseCalls, 1)
}
