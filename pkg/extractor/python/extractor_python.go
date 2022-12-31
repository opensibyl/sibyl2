package python

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

type Extractor struct {
}

type FunctionExtras struct {
	Decorators []string `json:"decorators"`
}

func (extractor *Extractor) GetLang() core.LangType {
	return core.LangPython
}
