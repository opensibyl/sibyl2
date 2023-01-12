package kotlin

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

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

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*object.Clazz, error) {
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
