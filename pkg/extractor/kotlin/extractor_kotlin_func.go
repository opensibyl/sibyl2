package kotlin

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

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
	funcUnit.Namespace = pkgName

	funcIdentifier := core.FindFirstByKindInSubsWithBfs(unit, KindKotlinSimpleIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func id found in identifier" + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content
	funcUnit.DefLine = int(funcIdentifier.Span.Start.Row + 1)

	return funcUnit, nil
}
