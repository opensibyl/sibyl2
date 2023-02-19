package java

import (
	"github.com/opensibyl/sibyl2/pkg/core"
)

// https://github.com/tree-sitter/tree-sitter-java/tree/master/src/node-types.json
const (
	KindJavaProgram              core.KindRepr = "program"
	KindJavaProgramDeclaration   core.KindRepr = "package_declaration"
	KindJavaScopeIdentifier      core.KindRepr = "scoped_identifier"
	KindJavaIdentifier           core.KindRepr = "identifier"
	KindJavaClassDeclaration     core.KindRepr = "class_declaration"
	KindJavaClassBody            core.KindRepr = "class_body"
	KindJavaFieldDeclaration     core.KindRepr = "field_declaration"
	KindJavaEnumDeclaration      core.KindRepr = "enum_declaration"
	KindJavaInterfaceDeclaration core.KindRepr = "interface_declaration"
	KindJavaMethodDeclaration    core.KindRepr = "method_declaration"
	KindJavaFormalParameters     core.KindRepr = "formal_parameters"
	KindJavaFormalParameter      core.KindRepr = "formal_parameter"
	KindJavaMethodInvocation     core.KindRepr = "method_invocation"
	KindJavaModifiers            core.KindRepr = "modifiers"
	KindJavaAnnotation           core.KindRepr = "annotation"
	KindJavaMarkerAnnotation     core.KindRepr = "marker_annotation"
	KindJavaBlock                core.KindRepr = "block"
	KindJavaSuperClass           core.KindRepr = "superclass"
	KindJavaSuperInterface       core.KindRepr = "super_interfaces"
	KindJavaTypeList             core.KindRepr = "type_list"
	KindJavaTypeIdentifier       core.KindRepr = "type_identifier"
	KindJavaGenericType          core.KindRepr = "generic_type"
	FieldJavaType                core.KindRepr = "type"
	FieldJavaDimensions          core.KindRepr = "dimensions"
	FieldJavaObject              core.KindRepr = "object"
	FieldJavaName                core.KindRepr = "name"
	FieldJavaArguments           core.KindRepr = "arguments"
	FieldJavaDeclarator          core.KindRepr = "declarator"
)

type Extractor struct {
}

func (extractor *Extractor) GetLang() core.LangType {
	return core.LangJava
}
