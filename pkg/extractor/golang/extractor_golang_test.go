package golang

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
	"github.com/stretchr/testify/assert"
)

var goCode = `
package abc

type Parser struct {
	*Headless
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

func injectV1Group(v1group *gin.RouterGroup) {
	// scope
	scopeGroup := v1group.Group("/")
}
`

func TestGolangExtractor_ExtractFunctions(t *testing.T) {
	t.Parallel()
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
	location := "8:46,10:1"
	if target.BodySpan.String() != location {
		panic(fmt.Sprintf("%s != %s", target.BodySpan.String(), location))
	}
	assert.Equal(t, target.Lang, core.LangGo)

	privateMethod := funcs[len(funcs)-1]
	assert.NotEqual(t, privateMethod.BodySpan.Start.Row, 0)
	assert.Equal(t, privateMethod.Namespace, "abc")
}

func TestGolangExtractor_Serialize(t *testing.T) {
	t.Parallel()
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
		normal, err := json.Marshal(each)
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
	t.Parallel()
	parser := core.NewParser(core.LangGo)
	units, err := parser.Parse([]byte(goCode))
	assert.Nil(t, err)

	extractor := &Extractor{}
	data, err := extractor.ExtractClasses(units)
	assert.Nil(t, err)
	for _, eachType := range data {
		core.Log.Infof("clazz: %v", eachType.GetSignature())

		fields := eachType.Extras.(*ClassExtras).Fields
		core.Log.Infof("fields: %v, %v", fields[0].Type, fields[0].Name)
	}
}
