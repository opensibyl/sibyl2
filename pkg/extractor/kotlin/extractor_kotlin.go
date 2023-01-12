package kotlin

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

// NOTICE: kotlin grammar is not official
// https://github.com/fwcd/tree-sitter-kotlin/blob/main/src/node-types.json
const (
	KindKotlinFunctionDecl   core.KindRepr = "function_declaration"
	KindKotlinFunctionBody   core.KindRepr = "function_body"
	KindKotlinPackageHeader  core.KindRepr = "package_header"
	KindKotlinIdentifier     core.KindRepr = "identifier"
	KindKotlinTypeIdentifier core.KindRepr = "type_identifier"
	KindKotlinClassDecl      core.KindRepr = "class_declaration"
	KindKotlinSourceFile     core.KindRepr = "source_file"
)

type Extractor struct {
}

func (extractor *Extractor) GetLang() core.LangType {
	return core.LangKotlin
}
