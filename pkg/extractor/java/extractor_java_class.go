package java

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type ClassField struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Annotations []string `json:"annotations"`
}

type ClassExtras struct {
	Annotations []string      `json:"annotations"`
	Fields      []*ClassField `json:"fields"`
}

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	if unit.Kind == KindJavaClassDeclaration || unit.Kind == KindJavaEnumDeclaration {
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

	program := core.FindFirstByKindInParent(unit, KindJavaProgram)
	packageDecl := core.FindFirstByKindInSubsWithDfs(program, KindJavaProgramDeclaration)
	packageIdentifier := core.FindFirstByKindInSubsWithDfs(packageDecl, KindJavaScopeIdentifier)
	if packageIdentifier == nil {
		core.Log.Warnf("no package found in %s", unit.Content)
	} else {
		clazz.Module = packageIdentifier.Content
	}

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindJavaClassDeclaration, KindJavaEnumDeclaration, KindJavaInterfaceDeclaration)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazz.Name = clazzIdentifier.Content

	extras := &ClassExtras{}
	// class annotations
	classModifiers := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaModifiers)
	if classModifiers != nil {
		classAnnotations := core.FindAllByKindsInSubs(classModifiers, KindJavaMarkerAnnotation, KindJavaAnnotation)
		if len(classAnnotations) != 0 {
			for _, each := range classAnnotations {
				extras.Annotations = append(extras.Annotations, each.Content)
			}
		}
	}
	// todo: fields
	clazz.Extras = extras

	return clazz, nil
}
