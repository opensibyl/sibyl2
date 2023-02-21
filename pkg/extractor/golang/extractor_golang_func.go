package golang

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
	"golang.org/x/exp/slices"
)

type FuncExtras struct {
}

func (extractor *Extractor) IsFunction(unit *core.Unit) bool {
	allowed := []core.KindRepr{
		KindGolangMethodDecl,
		KindGolangFuncDecl,
	}
	return slices.Contains(allowed, unit.Kind)
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
	switch unit.Kind {
	case KindGolangFuncDecl:
		return extractor.funcUnit2Function(unit)
	case KindGolangMethodDecl:
		return extractor.methodUnit2Function(unit)
	default:
		// should not reach here
		return nil, errors.New("IMPOSSIBLE")
	}
}

func (extractor *Extractor) methodUnit2Function(unit *core.Unit) (*object.Function, error) {
	funcUnit := &object.Function{}
	funcUnit.Span = unit.Span
	funcUnit.Unit = unit

	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindGolangBlock)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	// name
	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangFieldIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// receiver
	parameterList := core.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	parameterList = core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterList)
	receiverDecl := core.FindFirstByKindInSubsWithDfs(parameterList, KindGolangParameterDecl)
	typeDecl := core.FindFirstByFieldInSubsWithDfs(receiverDecl, FieldGolangType)
	if typeDecl == nil {
		return nil, errors.New("no receiver found in: " + typeDecl.Content)
	}
	funcUnit.Receiver = typeDecl.Content

	// namespace: package
	root := core.FindFirstByKindInParent(unit, KindGolangSourceFile)
	pkgName := core.FindFirstByKindInSubsWithDfs(root, KindGolangPackageIdentifier)
	if pkgName != nil {
		funcUnit.Namespace = pkgName.Content
	}

	// params
	paramListList := core.FindAllByKindInSubsWithDfs(unit, KindGolangParameterList)
	// no param == empty slice, never nil
	paramList := paramListList[1]
	for _, each := range core.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}

		valueUnit := &object.ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	// never nil
	retParams := core.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	switch retParams.Kind {
	case KindGolangParameterList:
		// multi params
		for _, each := range core.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
			typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
			paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
			var paramNameContent string
			if paramName == nil {
				paramNameContent = ""
			} else {
				paramNameContent = paramName.Content
			}
			valueUnit := &object.ValueUnit{
				Type: typeName.Content,
				Name: paramNameContent,
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		}
	case KindGolangTypeIdentifier:
		// only one param, and anonymous
		valueUnit := &object.ValueUnit{
			Type: retParams.Content,
			Name: "",
		}
		funcUnit.Returns = append(funcUnit.Returns, valueUnit)
	default:
		// no returns
	}
	// extras
	funcUnit.Extras = &FuncExtras{}

	return funcUnit, nil
}

func (extractor *Extractor) funcUnit2Function(unit *core.Unit) (*object.Function, error) {
	funcUnit := &object.Function{}
	funcUnit.Span = unit.Span
	funcUnit.Lang = extractor.GetLang()
	funcUnit.Unit = unit

	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindGolangBlock)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	// name
	funcIdentifier := core.FindFirstByKindInSubsWithDfs(unit, KindGolangIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func name found in " + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// namespace: package
	root := core.FindFirstByKindInParent(unit, KindGolangSourceFile)
	pkgName := core.FindFirstByKindInSubsWithDfs(root, KindGolangPackageIdentifier)
	if pkgName != nil {
		funcUnit.Namespace = pkgName.Content
	}

	// params
	// no param == empty slice, never nil
	paramList := core.FindFirstByKindInSubsWithDfs(unit, KindGolangParameterList)
	for _, each := range core.FindAllByKindInSubsWithDfs(paramList, KindGolangParameterDecl) {
		typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
		paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
		var paramNameContent string
		if paramName == nil {
			paramNameContent = ""
		} else {
			paramNameContent = paramName.Content
		}
		valueUnit := &object.ValueUnit{
			Type: typeName.Content,
			Name: paramNameContent,
		}
		funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
	}

	// returns
	retParams := core.FindFirstByFieldInSubsWithDfs(unit, FieldGolangParameters)
	if retParams != nil {
		switch retParams.Kind {
		case KindGolangParameterList:
			// multi params
			for _, each := range core.FindAllByKindInSubsWithDfs(retParams, KindGolangParameterDecl) {
				typeName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangType)
				paramName := core.FindFirstByFieldInSubsWithDfs(each, FieldGolangName)
				var paramNameContent string
				if paramName == nil {
					paramNameContent = ""
				} else {
					paramNameContent = paramName.Content
				}
				valueUnit := &object.ValueUnit{
					Type: typeName.Content,
					Name: paramNameContent,
				}
				funcUnit.Returns = append(funcUnit.Returns, valueUnit)
			}
		case KindGolangTypeIdentifier:
			// only one param, and anonymous
			valueUnit := &object.ValueUnit{
				Type: retParams.Content,
				Name: "",
			}
			funcUnit.Returns = append(funcUnit.Returns, valueUnit)
		default:
			// no returns
		}
	}

	// extras
	funcUnit.Extras = &FuncExtras{}

	return funcUnit, nil
}
