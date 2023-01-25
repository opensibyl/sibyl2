package javascript

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

// https://github.com/tree-sitter/tree-sitter-javascript/blob/master/src/node-types.json
const (
	KindJavaScriptClassDeclaration    core.KindRepr = "class_declaration"
	KindJavaScriptMethodDefinition    core.KindRepr = "method_definition"
	KindJavaScriptFunctionDeclaration core.KindRepr = "function_declaration"
	KindJavaScriptIdentifier          core.KindRepr = "identifier"
	KindJavaScriptFormalParameters    core.KindRepr = "formal_parameters"
	FieldJavaScriptName               core.KindRepr = "name"
	FieldJavaScriptParameters         core.KindRepr = "parameters"
)

type Extractor struct {
}

func (extractor *Extractor) GetLang() core.LangType {
	return core.LangJavaScript
}
