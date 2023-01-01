package golang

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type Field struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ClassExtras struct {
	Fields []*Field `json:"fields"`
}

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	// golang has no class. We use struct here.
	if unit.Kind == KindGolangTypeSpec {
		for _, eachSub := range unit.SubUnits {
			if eachSub.Kind == KindGolangStructType {
				return true
			}
		}
		return false
	}
	return false
}

func (extractor *Extractor) ExtractClasses(units []*core.Unit) ([]*object.Clazz, error) {
	var ret []*object.Clazz
	for _, eachUnit := range units {
		if !extractor.IsClass(eachUnit) {
			continue
		}
		eachClazz, err := extractor.unit2Clazz(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachClazz)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*object.Clazz, error) {
	data, err := extractor.ExtractClasses([]*core.Unit{unit})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("unit not a class: " + unit.Content)
	}
	return data[0], nil
}

func (extractor *Extractor) unit2Clazz(unit *core.Unit) (*object.Clazz, error) {
	clazz := object.NewClazz()

	// struct name
	structName := core.FindFirstByFieldInSubsWithBfs(
		core.FindFirstByKindInSubsWithBfs(unit, KindGolangTypeSpec),
		FieldGolangName)
	clazz.Name = structName.Content

	// package name
	root := core.FindFirstByKindInParent(unit, KindGolangSourceFile)
	pkgName := core.FindFirstByKindInSubsWithDfs(root, KindGolangPackageIdentifier)
	if pkgName != nil {
		clazz.Module = pkgName.Content
	}

	extras := &ClassExtras{}
	fieldList := core.FindFirstByKindInSubsWithBfs(unit, KindGolangFieldDeclList)

	if fieldList != nil {
		fields := core.FindAllByKindInSubs(fieldList, KindGolangFieldDecl)
		for _, eachField := range fields {
			typeDef := core.FindFirstByFieldInSubs(eachField, FieldGolangType)
			nameDef := core.FindFirstByFieldInSubs(eachField, FieldGolangName)
			f := &Field{}
			if typeDef != nil {
				f.Type = typeDef.Content
			}
			if nameDef != nil {
				f.Name = nameDef.Content
			}
			extras.Fields = append(extras.Fields, f)
		}
	}
	clazz.Extras = extras

	return clazz, nil
}
