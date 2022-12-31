package extractor

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

// https://github.com/tree-sitter/tree-sitter-python/blob/master/src/node-types.json
const (
	KindPythonFunctionDefinition  core.KindRepr = "function_definition"
	KindPythonIdentifier          core.KindRepr = "identifier"
	KindPythonDecoratedDefinition core.KindRepr = "decorated_definition"
	KindPythonDecorator           core.KindRepr = "decorator"
	KindPythonBlock               core.KindRepr = "block"
)

type PythonExtractor struct {
}

type PythonFunctionExtras struct {
	Decorators []string `json:"decorators"`
}

func (extractor *PythonExtractor) GetLang() core.LangType {
	return core.LangPython
}
