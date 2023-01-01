package golang

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

var goCode = `
package abc

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

	extractor := &Extractor{}
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

	extractor := &Extractor{}
	functions, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
	for _, each := range functions {
		normal, err := each.ToJson()
		if err != nil {
			panic(err)
		}

		var m map[string]any
		err = json.Unmarshal(normal, &m)
		if err != nil {
			panic(err)
		}
		back, err := object.Map2Func(m)
		if err != nil {
			panic(err)
		}
		core.Log.Infof("before func %v", each)
		core.Log.Infof("after func %v", back)
		if each.Name != back.Name {
			panic(errors.New("CONVERT FAILED"))
		}
	}
}

func TestExtractor_ExtractClasses(t *testing.T) {
	parser := core.NewParser(core.LangGo)
	units, err := parser.Parse([]byte(goCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	data, err := extractor.ExtractClasses(units)
	for _, eachType := range data {
		core.Log.Infof("clazz: %v", eachType.GetSignature())

		fields := eachType.Extras.(*ClassExtras).Fields
		core.Log.Infof("fields: %v, %v", fields[0].Type, fields[0].Name)
	}
}
