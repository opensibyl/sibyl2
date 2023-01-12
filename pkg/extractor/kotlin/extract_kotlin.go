package kotlin

import (
	"errors"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
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

func (extractor *Extractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractSymbols(units []*core.Unit) ([]*extractor.Symbol, error) {
	ret := make([]*object.Symbol, 0)
	for _, eachUnit := range units {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &object.Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
			Unit:      eachUnit,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
}

func (extractor *Extractor) IsFunction(unit *core.Unit) bool {
	if unit.Kind == KindKotlinFunctionDecl {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractFunctions(units []*core.Unit) ([]*object.Function, error) {
	var ret []*object.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc, err := extractor.ExtractFunction(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractFunction(unit *core.Unit) (*object.Function, error) {
	funcUnit := object.NewFunction()
	funcUnit.Span = unit.Span
	funcUnit.Unit = unit
	funcUnit.Lang = extractor.GetLang()

	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindKotlinFunctionBody)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	var pkgName string
	var clazzName string

	// trace its package
	root := core.FindFirstByKindInParent(unit, KindKotlinSourceFile)
	packageDecl := core.FindFirstByKindInSubsWithDfs(root, KindKotlinPackageHeader)
	packageIdentifier := core.FindFirstByKindInSubsWithDfs(packageDecl, KindKotlinIdentifier)
	if packageIdentifier == nil {
		core.Log.Warnf("no package found in %s", unit.Content)
	} else {
		pkgName = packageIdentifier.Content
	}

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindKotlinClassDecl)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindKotlinTypeIdentifier)
	if clazzIdentifier != nil {
		// in kotlin, function can decl without class
		clazzName = clazzIdentifier.Content
		funcUnit.Receiver = pkgName + "." + clazzName
	}

	return funcUnit, nil
}

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	// current kotlin grammar only has one class decl type (no interface and something else
	if unit.Kind == KindKotlinClassDecl {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractClasses(units []*core.Unit) ([]*object.Clazz, error) {
	var ret []*object.Clazz
	for _, eachUnit := range units {
		if !extractor.IsClass(eachUnit) {
			continue
		}
		eachClazz, err := extractor.ExtractClass(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachClazz)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*extractor.Clazz, error) {
	clazz := object.NewClazz()
	clazz.Span = unit.Span
	clazz.Lang = extractor.GetLang()
	clazz.Unit = unit

	var pkgName string
	var clazzName string

	// trace its package
	root := core.FindFirstByKindInParent(unit, KindKotlinSourceFile)
	packageDecl := core.FindFirstByKindInSubsWithDfs(root, KindKotlinPackageHeader)
	packageIdentifier := core.FindFirstByKindInSubsWithDfs(packageDecl, KindKotlinIdentifier)
	if packageIdentifier == nil {
		core.Log.Warnf("no package found in %s", unit.Content)
	} else {
		pkgName = packageIdentifier.Content
	}

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindKotlinClassDecl)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindKotlinTypeIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazzName = clazzIdentifier.Content

	clazz.Module = pkgName
	clazz.Name = clazzName
	return clazz, nil
}

func (extractor *Extractor) IsCall(unit *core.Unit) bool {
	//TODO implement me
	panic("implement me")
}

func (extractor *Extractor) ExtractCalls(units []*core.Unit) ([]*extractor.Call, error) {
	//TODO implement me
	panic("implement me")
}
