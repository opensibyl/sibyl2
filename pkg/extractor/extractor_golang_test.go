package extractor

import (
	"errors"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
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
	funcs, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}

	// check its base info
	target := funcs[0]
	core.Log.Debugf("target: %s, span: %s", target.Name, target.BodySpan.String())
	if target.BodySpan.String() != "5:46,7:1" {
		panic(nil)
	}
}

func TestGolangExtractor_Serialize(t *testing.T) {
	parser := core.NewParser(core.LangGo)
	units, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.LangGo)
	functions, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
	for _, each := range functions {
		normal, err := each.ToJson()
		if err != nil {
			panic(err)
		}

		back, err := FromJson(normal)
		if err != nil {
			panic(err)
		}
		core.Log.Debugf("before func %v", each)
		core.Log.Debugf("after func %v", back)
		if each.Name != back.Name {
			panic(errors.New("CONVERT FAILED"))
		}
	}
}
