package java

import (
	"errors"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type ClassField struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Annotations []string `json:"annotations"`
	Modifiers   []string `json:"modifiers"`
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
	clazz.Lang = extractor.GetLang()
	clazz.Unit = unit

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
	// fields
	body := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindJavaClassBody)
	if body != nil {
		fields := core.FindAllByKindsInSubs(body, KindJavaFieldDeclaration)
		for _, eachField := range fields {
			typeDecl := core.FindFirstByFieldInSubs(eachField, FieldJavaType)
			variableDecl := core.FindFirstByFieldInSubs(eachField, FieldJavaDeclarator)
			nameDecl := core.FindFirstByKindInSubsWithBfs(variableDecl, KindJavaIdentifier)
			if nameDecl == nil || typeDecl == nil {
				return nil, errors.New("not finished field decl")
			}
			field := &ClassField{
				Name:        nameDecl.Content,
				Type:        typeDecl.Content,
				Annotations: nil,
				Modifiers:   nil,
			}
			extras.Fields = append(extras.Fields, field)

			modifiers := core.FindFirstByKindInSubsWithBfs(eachField, KindJavaModifiers)
			if modifiers == nil {
				// no modifiers and annotations
				continue
			}
			modifiersStr := modifiers.Content

			// annotation?
			annotations := core.FindAllByKindsInSubs(modifiers, KindJavaMarkerAnnotation, KindJavaAnnotation)
			if len(annotations) != 0 {
				for _, each := range annotations {
					field.Annotations = append(extras.Annotations, each.Content)
					// remove it from modifiers
					// currently tree-sitter did not split these nodes
					modifiersStr = strings.Replace(modifiersStr, each.Content, "", 1)
				}
			}
			field.Modifiers = strings.Split(strings.TrimSpace(modifiersStr), " ")
		}
	}
	clazz.Extras = extras

	return clazz, nil
}
