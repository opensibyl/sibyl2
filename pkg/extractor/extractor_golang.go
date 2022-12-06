package extractor

import (
	"github.com/williamfzc/sibyl2/pkg/core"
)

// https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
const (
	KindGolangMethodDecl      core.KindRepr = "method_declaration"
	KindGolangFuncDecl        core.KindRepr = "function_declaration"
	KindGolangIdentifier      core.KindRepr = "identifier"
	KindGolangFieldIdentifier core.KindRepr = "field_identifier"
	KindGolangTypeIdentifier  core.KindRepr = "type_identifier"
	KindGolangParameterList   core.KindRepr = "parameter_list"
	KindGolangParameterDecl   core.KindRepr = "parameter_declaration"
	KindGolangCallExpression  core.KindRepr = "call_expression"
	FieldGolangType           core.KindRepr = "type"
	FieldGolangName           core.KindRepr = "name"
	FieldGolangParameters     core.KindRepr = "parameters"
	FieldGolangFunction       core.KindRepr = "function"
	FieldGolangArguments      core.KindRepr = "arguments"
	FieldGolangResult         core.KindRepr = "result"
)

type GolangExtractor struct {
}

func (extractor *GolangExtractor) GetLang() core.LangType {
	return core.LangGo
}
