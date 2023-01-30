package golang

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

// https://github.com/tree-sitter/tree-sitter-go/blob/master/src/node-types.json
const (
	KindGolangMethodDecl        core.KindRepr = "method_declaration"
	KindGolangFuncDecl          core.KindRepr = "function_declaration"
	KindGolangIdentifier        core.KindRepr = "identifier"
	KindGolangFieldIdentifier   core.KindRepr = "field_identifier"
	KindGolangTypeIdentifier    core.KindRepr = "type_identifier"
	KindGolangParameterList     core.KindRepr = "parameter_list"
	KindGolangParameterDecl     core.KindRepr = "parameter_declaration"
	KindGolangCallExpression    core.KindRepr = "call_expression"
	KindGolangTypeSpec          core.KindRepr = "type_spec"
	KindGolangStructType        core.KindRepr = "struct_type"
	KindGolangFieldDeclList     core.KindRepr = "field_declaration_list"
	KindGolangFieldDecl         core.KindRepr = "field_declaration"
	KindGolangPackageIdentifier core.KindRepr = "package_identifier"
	KindGolangSourceFile        core.KindRepr = "source_file"
	KindGolangBlock             core.KindRepr = "block"
	FieldGolangType             core.KindRepr = "type"
	FieldGolangName             core.KindRepr = "name"
	FieldGolangParameters       core.KindRepr = "parameters"
	FieldGolangFunction         core.KindRepr = "function"
	FieldGolangArguments        core.KindRepr = "arguments"
)

type Extractor struct {
}

func (extractor *Extractor) GetLang() core.LangType {
	return core.LangGo
}
