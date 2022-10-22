package extractor

import (
	"fmt"
	"github.com/williamfzc/sibyl2/pkg/core"
	"testing"
)

var pythonCode = `
import requests

def a():
	b("abcde")

@DDDDeco
def b(s):
	print("defabc")

class C(object):
	pass
`

func TestPythonExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.LangPython)
	units, err := parser.Parse([]byte(pythonCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.LangPython)
	functions, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
	for _, each := range functions {
		fmt.Printf("%+v", each.Extras)
	}
}
